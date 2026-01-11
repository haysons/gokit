package transport

import (
	"context"
)

// Server 传输层 server 部分
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Header 传输层 header 部分，主要用于存取元数据
type Header interface {
	Get(key string) string
	Set(key string, value string)
	Add(key string, value string)
	Keys() []string
	Values(key string) []string
}

type Transporter interface {
	// Kind http 或 grpc
	Kind() Kind

	// Operation 返回当前的服务名及方法名，可直观定位当前操作
	Operation() string

	// RequestHeader 返回请求 header
	// http: http.Header
	// grpc: metadata.MD
	RequestHeader() Header

	// ReplyHeader 返回响应 header
	// http: http.Header
	// grpc: metadata.MD
	ReplyHeader() Header
}

// Kind 传输层类型，grpc 或 http
type Kind string

func (k Kind) String() string { return string(k) }

const (
	KindGRPC Kind = "grpc"
	KindHTTP Kind = "http"
)

type (
	serverTransportKey struct{}
	clientTransportKey struct{}
)

// InjectServerContext 向 ctx 中注入传输层
func InjectServerContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, serverTransportKey{}, tr)
}

// FromServerContext 自 context 中获取传输层
func FromServerContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(serverTransportKey{}).(Transporter)
	return
}

// InjectClientContext 向 ctx 中注入传输层
func InjectClientContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, clientTransportKey{}, tr)
}

// FromClientContext 自 context 中获取传输层
func FromClientContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(clientTransportKey{}).(Transporter)
	return
}
