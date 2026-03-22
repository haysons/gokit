package http

import (
	"context"

	"github.com/haysons/gokit/transport/testdata/helloworld"
)

// =============================================================================
// Middleware Support for HTTP Transport
// =============================================================================

// Interceptor is a function that wraps the handler to add pre/post processing.
type Interceptor func(ctx context.Context, methodName string, req interface{}, handler func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error)

// globalInterceptors holds the registered interceptors
var globalInterceptors []Interceptor

// RegisterInterceptor registers a global interceptor (first registered = outermost)
func RegisterInterceptor(i Interceptor) {
	globalInterceptors = append(globalInterceptors, i)
}

// ClearInterceptors clears all registered interceptors
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

// =============================================================================
// GreeterServerWrapper - Generated wrapper for Greeter service
// =============================================================================

// GreeterServerWrapper wraps a GreeterServer to add interceptor support.
type GreeterServerWrapper struct {
	helloworld.UnimplementedGreeterServer
	impl helloworld.GreeterServer
}

// NewGreeterServerWrapper creates a GreeterServer wrapper with interceptors.
func NewGreeterServerWrapper(impl helloworld.GreeterServer) *GreeterServerWrapper {
	return &GreeterServerWrapper{impl: impl}
}

// Use adds interceptors to the wrapper.
func (w *GreeterServerWrapper) Use(interceptors ...Interceptor) {
	for _, i := range interceptors {
		RegisterInterceptor(i)
	}
}

// SayHello implements helloworld.GreeterServer with interceptor support.
func (w *GreeterServerWrapper) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	result, err := callWithInterceptors(ctx, "SayHello", req, func(ctx context.Context, req interface{}) (interface{}, error) {
		return w.impl.SayHello(ctx, req.(*helloworld.HelloRequest))
	})
	if err != nil {
		return nil, err
	}
	return result.(*helloworld.HelloReply), nil
}

// =============================================================================
// ServerWrapper - Generic wrapper generator
// =============================================================================

// ServerWrapper wraps any gRPC server to add interceptor support.
// This is a helper that can be used with code generation.
//
// For each service, create a wrapper like this:
//
//   type XxxServerWrapper struct {
//       helloworld.UnimplementedXxxServer
//       impl helloworld.XxxServer
//   }
//
//   func (w *XxxServerWrapper) YyyMethod(ctx context.Context, req *YyyRequest) (*YyyResponse, error) {
//       result, err := callWithInterceptors(ctx, "YyyMethod", req, func(ctx, req) {
//           return w.impl.YyyMethod(ctx, req.(*YyyRequest))
//       })
//       if err != nil {
//           return nil, err
//       }
//       return result.(*YyyResponse), nil
//   }
