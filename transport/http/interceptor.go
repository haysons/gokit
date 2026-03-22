package http

import (
	"context"
)

// =============================================================================
// Interceptor - Middleware for HTTP Gateway
// =============================================================================

// Interceptor is a function that wraps the handler to add pre/post processing.
// It's called after the request is decoded but before the business method is invoked.
type Interceptor func(ctx context.Context, methodName string, req interface{}, handler func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error)

// globalInterceptors holds the registered interceptors
var globalInterceptors []Interceptor

// RegisterInterceptor registers a global interceptor that will be called for all methods.
// Interceptors are executed in reverse registration order (last registered = outermost).
func RegisterInterceptor(i Interceptor) {
	globalInterceptors = append(globalInterceptors, i)
}

// ClearInterceptors clears all registered interceptors (useful for testing)
func ClearInterceptors() {
	globalInterceptors = nil
}

// GetInterceptors returns the registered interceptors
func GetInterceptors() []Interceptor {
	return globalInterceptors
}

// callWithInterceptors calls the interceptor chain for a method.
func callWithInterceptors(ctx context.Context, methodName string, req interface{}, handler func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error) {
	h := handler
	for i := len(globalInterceptors) - 1; i >= 0; i-- {
		next := h
		m := globalInterceptors[i]
		h = func(ctx context.Context, req interface{}) (interface{}, error) {
			return m(ctx, methodName, req, func(ctx context.Context, req interface{}) (interface{}, error) {
				return next(ctx, req)
			})
		}
	}
	return h(ctx, req)
}
