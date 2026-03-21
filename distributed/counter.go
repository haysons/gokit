package distributed

import (
	"context"
	"fmt"
	"strconv"

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
	// 先读取当前值
	getResp, err := c.client.Get(ctx, c.key)
	if err != nil {
		return 0, fmt.Errorf("incr counter failed to get current value: %w", err)
	}

	var prevVal int64
	var modRevision int64

	if getResp.Count > 0 {
		prevVal, err = strconv.ParseInt(string(getResp.Kvs[0].Value), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("incr counter failed to parse current value: %w", err)
		}
		modRevision = getResp.Kvs[0].ModRevision
	}

	newVal := prevVal + delta
	newValStr := strconv.FormatInt(newVal, 10)

	// 使用 CAS 事务保证原子性：若 modRevision 未变则写入新值，否则说明有并发修改
	var txnResp *clientv3.TxnResponse
	if modRevision == 0 {
		// key 不存在，条件：version == 0
		txnResp, err = c.client.Txn(ctx).
			If(clientv3.Compare(clientv3.Version(c.key), "=", 0)).
			Then(clientv3.OpPut(c.key, newValStr)).
			Commit()
	} else {
		// key 存在，条件：modRevision 未变
		txnResp, err = c.client.Txn(ctx).
			If(clientv3.Compare(clientv3.ModRevision(c.key), "=", modRevision)).
			Then(clientv3.OpPut(c.key, newValStr)).
			Commit()
	}
	if err != nil {
		return 0, fmt.Errorf("incr counter failed: %w", err)
	}

	if !txnResp.Succeeded {
		// 发生并发竞争，重试
		return c.IncrBy(ctx, delta)
	}

	return newVal, nil
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
