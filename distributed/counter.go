package distributed

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.etcd.io/etcd/client/v3"
)

// Counter 分布式计数器
type Counter struct {
	client *clientv3.Client
	key    string
}

// NewCounter 创建分布式计数器实例
func NewCounter(client *clientv3.Client, key string) *Counter {
	return &Counter{
		client: client,
		key:    key,
	}
}

// Get 获取当前计数值
func (c *Counter) Get(ctx context.Context) (int64, error) {
	resp, err := c.client.Get(ctx, c.key)
	if err != nil {
		return 0, fmt.Errorf("get counter failed: %w", err)
	}

	if resp.Count == 0 {
		return 0, nil
	}

	value, err := strconv.ParseInt(string(resp.Kvs[0].Value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse counter value failed: %w", err)
	}

	return value, nil
}

// Incr 原子递增计数，返回新值
func (c *Counter) Incr(ctx context.Context) (int64, error) {
	return c.IncrBy(ctx, 1)
}

// IncrBy 原子递增指定值，返回新值
func (c *Counter) IncrBy(ctx context.Context, delta int64) (int64, error) {
	op := clientv3.OpTxn().
		If(clientv3.Compare(clientv3.Version(c.key), ">", 0)).
		Then(clientv3.OpPut(c.key, strconv.FormatInt(delta, 10), clientv3.WithPrevKV())).
		Else(clientv3.OpPut(c.key, strconv.FormatInt(delta, 10)))

	resp, err := c.client.Txn(ctx).Commit(op)
	if err != nil {
		return 0, fmt.Errorf("incr counter failed: %w", err)
	}

	if resp.Succeeded {
		prevVal, _ := strconv.ParseInt(string(resp.Responses[0].GetResponsePut().PrevKv.Value), 10, 64)
		return prevVal + delta, nil
	}

	return delta, nil
}

// Decr 原子递减计数，返回新值
func (c *Counter) Decr(ctx context.Context) (int64, error) {
	return c.DecrBy(ctx, 1)
}

// DecrBy 原子递减指定值，返回新值
func (c *Counter) DecrBy(ctx context.Context, delta int64) (int64, error) {
	return c.IncrBy(ctx, -delta)
}

// Set 直接设置计数值
func (c *Counter) Set(ctx context.Context, value int64) error {
	_, err := c.client.Put(ctx, c.key, strconv.FormatInt(value, 10))
	if err != nil {
		return fmt.Errorf("set counter failed: %w", err)
	}

	return nil
}

// Reset 重置计数器到 0
func (c *Counter) Reset(ctx context.Context) error {
	_, err := c.client.Delete(ctx, c.key)
	if err != nil {
		return fmt.Errorf("reset counter failed: %w", err)
	}

	return nil
}
