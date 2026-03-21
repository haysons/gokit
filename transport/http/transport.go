package http

import (
	"net/http"

	"github.com/haysons/gokit/transport"
)

var _ transport.Transporter = (*Transport)(nil)

// Transport HTTP 传输层实现
type Transport struct {
	endpoint      string
	operation     string
	reqHeader     http.Header
	replyHeader   http.Header
	request       *http.Request
	response      http.ResponseWriter
	pathTemplate  string
}

// Kind 传输类别
func (tr *Transport) Kind() transport.Kind {
	return transport.KindHTTP
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
	return headerCarrier{tr.reqHeader}
}

// ReplyHeader 响应头
func (tr *Transport) ReplyHeader() transport.Header {
	return headerCarrier{tr.replyHeader}
}

// Request 返回原始 HTTP 请求
func (tr *Transport) Request() interface{} {
	return tr.request
}

// Response 返回原始 HTTP 响应
func (tr *Transport) Response() http.ResponseWriter {
	return tr.response
}

// PathTemplate returns the route path template
func (tr *Transport) PathTemplate() string {
	return tr.pathTemplate
}

// headerCarrier HTTP header 包装
type headerCarrier struct {
	http.Header
}

func (hc headerCarrier) Get(key string) string {
	return hc.Header.Get(key)
}

func (hc headerCarrier) Set(key, value string) {
	hc.Header.Set(key, value)
}

func (hc headerCarrier) Add(key, value string) {
	hc.Header.Add(key, value)
}

func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc.Header))
	for k := range hc.Header {
		keys = append(keys, k)
	}
	return keys
}

func (hc headerCarrier) Values(key string) []string {
	return hc.Header[key]
}
