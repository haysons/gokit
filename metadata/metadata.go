package metadata

import (
	"context"
	"fmt"
	"strings"
)

// Metadata 跨 rpc 传输的元数据
type Metadata map[string][]string

// New 基于 map 创建元数据
func New(mds ...map[string][]string) Metadata {
	md := Metadata{}
	for _, m := range mds {
		for k, vList := range m {
			for _, v := range vList {
				md.Add(k, v)
			}
		}
	}
	return md
}

// Add 元数据中增加 kv
func (m Metadata) Add(key, value string) {
	if key == "" {
		return
	}

	lowerKey := strings.ToLower(key)
	m[lowerKey] = append(m[lowerKey], value)
}

// Get 通过 key 获取元数据
func (m Metadata) Get(key string) string {
	v := m[strings.ToLower(key)]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Set 元数据中保存 kv
func (m Metadata) Set(key string, value string) {
	if key == "" || value == "" {
		return
	}
	m[strings.ToLower(key)] = []string{value}
}

// Range 遍历元数据中的全部 kv
func (m Metadata) Range(f func(k string, v []string) bool) {
	for k, v := range m {
		if !f(k, v) {
			break
		}
	}
}

// Values 获取元数据中 key 的全部值
func (m Metadata) Values(key string) []string {
	return m[strings.ToLower(key)]
}

// Clone 对元数据进行复制
func (m Metadata) Clone() Metadata {
	md := make(Metadata, len(m))
	for k, v := range m {
		md[k] = v
	}
	return md
}

type serverMetadataKey struct{}

// InjectServerContext 将元数据注入 server context 之中
func InjectServerContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, serverMetadataKey{}, md)
}

// FromServerContext 自 server context 之中提取元数据
func FromServerContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(serverMetadataKey{}).(Metadata)
	return md, ok
}

type clientMetadataKey struct{}

// InjectClientContext 将元数据注入 client context 之中
func InjectClientContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, clientMetadataKey{}, md)
}

// FromClientContext 自 client context 之中提取元数据
func FromClientContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(clientMetadataKey{}).(Metadata)
	return md, ok
}

// AppendToClientContext 在 client context 中添加 kv 列表
func AppendToClientContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToClientContext got an odd number of input pairs for metadata: %d", len(kv)))
	}
	md, _ := FromClientContext(ctx)
	md = md.Clone()
	for i := 0; i < len(kv); i += 2 {
		md.Set(kv[i], kv[i+1])
	}
	return InjectClientContext(ctx, md)
}

// MergeToClientContext 将元数据合并至 client context 之中
func MergeToClientContext(ctx context.Context, cmd Metadata) context.Context {
	md, _ := FromClientContext(ctx)
	md = md.Clone()
	for k, v := range cmd {
		md[k] = v
	}
	return InjectClientContext(ctx, md)
}
