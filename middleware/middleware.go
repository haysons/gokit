package middleware

import (
	"context"
)

// Handler 表示请求处理函数，接收 ctx 和请求参数，返回响应和错误。
type Handler func(ctx context.Context, req any) (any, error)

// Middleware 表示中间件，它接收下一级 Handler，并返回一个新的 Handler。
// 中间件通常在调用前后增加额外逻辑（如日志、鉴权、限流等）。
type Middleware func(next Handler) Handler

// Combine 将多个中间件合并为一个中间件。
// 中间件会按照传入顺序依次组合，执行时遵循“洋葱模型”——
// 最先传入的中间件最外层执行，最后传入的中间件最内层执行。
func Combine(m ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
