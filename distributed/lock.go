package distributed

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// Lock 分布式锁
type Lock struct {
	client     *clientv3.Client
	session    *concurrency.Session
	mutex      *concurrency.Mutex
	lockKey    string
}

// NewLock 创建分布式锁实例
func NewLock(client *clientv3.Client, lockKey string) (*Lock, error) {
	session, err := concurrency.NewSession(client, concurrency.WithTTL(60))
	if err != nil {
		return nil, fmt.Errorf("create session failed: %w", err)
	}

	return &Lock{
		client:     client,
		session:    session,
		mutex:      concurrency.NewMutex(session, lockKey),
		lockKey:    lockKey,
	}, nil
}

// Lock 尝试获取分布式锁，支持超时
func (l *Lock) Lock(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return l.mutex.Lock(ctx)
}

// TryLock 尝试获取分布式锁，立即返回结果
func (l *Lock) TryLock(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err := l.mutex.Lock(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// Unlock 释放分布式锁
func (l *Lock) Unlock(ctx context.Context) error {
	return l.mutex.Unlock(ctx)
}

// Close 关闭锁实例，清理资源
func (l *Lock) Close() error {
	return l.session.Close()
}
