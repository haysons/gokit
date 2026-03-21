package distributed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/client/v3"
)

// Queue 分布式队列
type Queue struct {
	client    *clientv3.Client
	keyPrefix string
}

// NewQueue 创建分布式队列实例
func NewQueue(client *clientv3.Client, keyPrefix string) *Queue {
	return &Queue{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// Enqueue 入队操作
func (q *Queue) Enqueue(ctx context.Context, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value failed: %w", err)
	}

	_, err = q.client.Put(ctx, fmt.Sprintf("%s/%d", q.keyPrefix, time.Now().UnixNano()), string(data))
	if err != nil {
		return fmt.Errorf("enqueue failed: %w", err)
	}

	return nil
}

// Dequeue 出队操作，支持超时
func (q *Queue) Dequeue(ctx context.Context, timeout time.Duration) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		// 获取队列中最旧的元素（按 key 升序取第一条）
		resp, err := q.client.Get(ctx, q.keyPrefix,
			clientv3.WithPrefix(),
			clientv3.WithLimit(1),
			clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		)
		if err != nil {
			return nil, fmt.Errorf("get queue elements failed: %w", err)
		}

		if resp.Count > 0 {
			// 尝试删除该元素（CAS：version 未变才删除，防止并发重复消费）
			key := string(resp.Kvs[0].Key)
			txnResp, err := q.client.Txn(ctx).
				If(clientv3.Compare(clientv3.Version(key), "=", resp.Kvs[0].Version)).
				Then(clientv3.OpDelete(key)).
				Commit()
			if err != nil {
				return nil, fmt.Errorf("dequeue transaction failed: %w", err)
			}

			if txnResp.Succeeded {
				var value interface{}
				if err := json.Unmarshal(resp.Kvs[0].Value, &value); err != nil {
					return nil, fmt.Errorf("unmarshal value failed: %w", err)
				}

				return value, nil
			}
		}

		// 队列为空，等待一段时间后重试
		select {
		case <-time.After(100 * time.Millisecond):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Length 获取队列长度
func (q *Queue) Length(ctx context.Context) (int64, error) {
	resp, err := q.client.Get(ctx, q.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return 0, fmt.Errorf("get queue length failed: %w", err)
	}

	return resp.Count, nil
}

// Clear 清空队列
func (q *Queue) Clear(ctx context.Context) error {
	_, err := q.client.Delete(ctx, q.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("clear queue failed: %w", err)
	}

	return nil
}
