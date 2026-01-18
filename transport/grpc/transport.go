package grpc

import (
	"github.com/haysons/gokit/transport"
	grpcmd "google.golang.org/grpc/metadata"
)

var _ transport.Transporter = (*Transport)(nil)

// Transport grpc 传输层
type Transport struct {
	endpoint    string
	operation   string
	reqHeader   headerCarrier
	replyHeader headerCarrier
}

// Kind 传输类别
func (tr *Transport) Kind() transport.Kind {
	return transport.KindGRPC
}

// Endpoint 访问的端点
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation 访问的方法
func (tr *Transport) Operation() string {
	return tr.operation
}

// RequestHeader 请求头
func (tr *Transport) RequestHeader() transport.Header {
	return tr.reqHeader
}

// ReplyHeader 响应头
func (tr *Transport) ReplyHeader() transport.Header {
	return tr.replyHeader
}

// headerCarrier grpc 传输层使用 metadata.MD 传输 header
type headerCarrier grpcmd.MD

// Get 自 header 中获取特定 key 的值
func (mc headerCarrier) Get(key string) string {
	vals := grpcmd.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// Set 在 header 中设置键值
func (mc headerCarrier) Set(key string, value string) {
	grpcmd.MD(mc).Set(key, value)
}

// Add 在 header 中添加键值
func (mc headerCarrier) Add(key string, value string) {
	grpcmd.MD(mc).Append(key, value)
}

// Keys 获取 header 中的全部 key 列表
func (mc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range grpcmd.MD(mc) {
		keys = append(keys, k)
	}
	return keys
}

// Values 获取 header 中的全部值列表
func (mc headerCarrier) Values(key string) []string {
	return grpcmd.MD(mc).Get(key)
}
