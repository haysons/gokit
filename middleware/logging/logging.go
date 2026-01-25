package logging

import (
	"context"
	"log/slog"
	"time"

	"github.com/haysons/gokit/middleware"
	"github.com/haysons/gokit/transport"
)

// Server is an server logging middleware.
func Server(logger slog.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			var (
				code      int32
				kind      string
				operation string
			)

			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			logger.Info(
				"api log",
				"kind", kind,
				"operation", operation,
				"code", code,
				"duration", time.Since(startTime),
			)
			return
		}
	}
}
