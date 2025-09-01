package recovery

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/haysons/gokit/middleware"
)

func Recovery() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (resp any, err error) {
			defer func() {
				if e := recover(); e != nil {
					slog.Error("panic recovered", "error", e, "stack", string(debug.Stack()))
					err = fmt.Errorf("panic recovered: %v", e)
				}
			}()
			return next(ctx, req)
		}
	}
}
