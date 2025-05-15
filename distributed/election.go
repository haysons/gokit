package distributed

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sync"
	"sync/atomic"
	"time"
)

// Election 是伴随应用程序始终的选举对象，针对于etcd选举对象进行了更多的异常处理，在某次选举出现异常时，会发起新的选举
type Election struct {
	id        string // 竞选者id，当前使用ip-uid作为唯一标识
	elKey     string
	isLeader  int32
	electCh   chan error
	readyCh   chan struct{}
	readyOnce sync.Once
	ctx       context.Context
	cancel    context.CancelFunc
	client    *clientv3.Client
	session   *concurrency.Session
}

// NewElection 新建选举对象
func NewElection(client *clientv3.Client, elKey string) (*Election, error) {
	id := fmt.Sprintf("%s-%s", "", xid.New().String())
	ctx, cancel := context.WithCancel(context.Background())
	return &Election{
		id:       id,
		elKey:    elKey,
		isLeader: 0,
		electCh:  make(chan error, 3),
		readyCh:  make(chan struct{}),
		ctx:      ctx,
		cancel:   cancel,
		client:   client,
	}, nil
}

// ID 当前竞选者的id
func (el *Election) ID() string {
	return el.id
}

// IsLeader 当前竞选者是否为leader
func (el *Election) IsLeader() bool {
	return atomic.LoadInt32(&el.isLeader) == 1
}

// Start 发起竞选，并等待出现leader，若等待时间过长将会返回错误
func (el *Election) Start(ctx context.Context) error {
	go el.run(ctx)
	return el.waitForReady()
}

func (el *Election) run(ctx context.Context) {
	for {
		select {
		case <-el.ctx.Done():
			return
		default:
			elect, electRes, err := el.elect()
			if err != nil {
				time.Sleep(5 * time.Second)
				break
			}
			if err = el.listen(ctx, elect, electRes); err != nil {
			}
			time.Sleep(time.Second)
		}
	}
}

// elect 发起竞选
func (el *Election) elect() (*concurrency.Election, chan error, error) {
	session, err := concurrency.NewSession(el.client, concurrency.WithTTL(10), concurrency.WithContext(el.ctx))
	if err != nil {
		return nil, nil, err
	}
	el.session = session
	electRes := make(chan error, 1)
	election := concurrency.NewElection(session, el.elKey)
	go func() {
		electRes <- election.Campaign(el.ctx, el.id)
	}()
	return election, electRes, nil
}

// listen 监听本次竞选的全部事件，直到监听出现异常才会停止阻塞，返回错误
func (el *Election) listen(ctx context.Context, elect *concurrency.Election, electRes chan error) error {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	leaderChan := elect.Observe(el.ctx)
	for {
		select {
		case err := <-electRes:
			// 本次选举发生问题，返回错误，发起新的选举
			if err != nil {
				return err
			} else {
			}
		case resp, ok := <-leaderChan:
			if !ok {
				return errors.New("elect observe chan closed")
			}
			if len(resp.Kvs) > 0 {
				curLeader := string(resp.Kvs[0].Value)
				el.setLeader(ctx, curLeader)
				el.readyOnce.Do(func() {
					// 首次监听到leader发生变化，则初始化完成
					close(el.readyCh)
				})
			}
		case <-ticker.C:
			// 此处为兜底逻辑，若监听leader变化的管道发生网络问题，可通过定期的查询判断出leader是否发生变化
			resp, err := elect.Leader(el.ctx)
			if err != nil {
				if errors.Is(err, concurrency.ErrElectionNoLeader) {
					// 不存在leader，发起新的选举
					return concurrency.ErrElectionNoLeader
				}
				continue
			}
			if resp != nil && len(resp.Kvs) > 0 {
				el.setLeader(ctx, string(resp.Kvs[0].Value))
			}
		case <-el.ctx.Done():
			return el.ctx.Err()
		}
	}
}

func (el *Election) setLeader(_ context.Context, curLeader string) {
	if el.id == curLeader {
		atomic.CompareAndSwapInt32(&el.isLeader, 0, 1)
	} else {
		atomic.CompareAndSwapInt32(&el.isLeader, 1, 0)
	}
}

// waitForReady 实际等待
func (el *Election) waitForReady() error {
	ctx, cancel := context.WithTimeout(el.ctx, 10*time.Second)
	defer cancel()
	select {
	case <-el.readyCh:
		return nil
	case <-ctx.Done():
		if err := el.Close(); err != nil {
			return err
		}
		return errors.New("wait for election ready timeout")
	}
}

// Close 关闭竞选，清理资源，让leader立刻结束任期，需在程序退出时调用，否则间隔5s才能重新选出leader
func (el *Election) Close() error {
	el.cancel()
	if el.session != nil {
		return el.session.Close()
	}
	return nil
}
