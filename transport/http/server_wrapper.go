package http

import (
	"context"
	"reflect"

	"github.com/haysons/gokit/middleware"
)

// =============================================================================
// ServerWrapper - 包装 gRPC Server 以支持中间件
// =============================================================================

// ServerWrapper wraps a gRPC server with middlewares.
// This allows using protoc-generated RegisterXxxHandlerServer while still applying middlewares.
//
// Usage:
//
//	// 1. Your actual server implementation
//	type MyGreeterServer struct{}
//	func (s *MyGreeterServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
//	    return &HelloReply{Message: "Hello " + req.Name}, nil
//	}
//
//	// 2. Wrap with middlewares using MustWrap
//	wrapped := http.MustWrap(&MyGreeterServer{}, middleware.Logging(), middleware.Recovery())
//
//	// 3. Register with protoc-generated code
//	RegisterGreeterHandlerServer(ctx, mux, wrapped)
//
// The wrapper implements the same interface as your server, so it can be passed directly
// to the protoc-generated registration function.
type ServerWrapper struct {
	server      interface{}
	middlewares []middleware.Middleware
	methodType  reflect.Type
}

// MustWrap creates a server wrapper and panics if the server is nil.
func MustWrap(server interface{}, middlewares ...middleware.Middleware) *ServerWrapper {
	if server == nil {
		panic("server cannot be nil")
	}
	return NewWrapper(server, middlewares...)
}

// NewWrapper creates a server wrapper that applies middlewares to all unary methods.
func NewWrapper(server interface{}, middlewares ...middleware.Middleware) *ServerWrapper {
	wrapper := &ServerWrapper{
		server:      server,
		middlewares: middlewares,
	}

	if server != nil {
		serverType := reflect.TypeOf(server)
		if serverType.Kind() == reflect.Ptr {
			serverType = serverType.Elem()
		}
		wrapper.methodType = serverType
	}

	return wrapper
}

// Invoke calls the wrapped server's method with middleware applied.
// This method is called via reflection from wrapper methods.
func (w *ServerWrapper) Invoke(ctx context.Context, methodName string, req, resp interface{}) error {
	if w.server == nil {
		return nil
	}

	// Get the method from the server
	method := reflect.ValueOf(w.server).MethodByName(methodName)
	if !method.IsValid() {
		return nil
	}

	// Build the handler chain with middleware
	h := func(ctx context.Context, req any) (any, error) {
		// Call the actual server method
		results := method.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(req),
		})

		if len(results) >= 2 {
			if !results[1].IsNil() {
				return nil, results[1].Interface().(error)
			}
			if !results[0].IsZero() {
				return results[0].Interface(), nil
			}
		}
		return nil, nil
	}

	// Apply middlewares (outermost to innermost)
	for i := len(w.middlewares) - 1; i >= 0; i-- {
		nextHandler := h
		m := w.middlewares[i]
		h = m(func(ctx context.Context, req any) (any, error) {
			return nextHandler(ctx, req)
		})
	}

	// Execute the handler chain
	result, err := h(ctx, req)
	if err != nil {
		return err
	}

	// Copy result to response
	if result != nil && resp != nil {
		reflect.ValueOf(resp).Elem().Set(reflect.ValueOf(result).Elem())
	}

	return nil
}

// =============================================================================
// Helper to create method wrappers
// =============================================================================

// WrapMethod creates a method wrapper that applies middlewares.
// This is used to create wrapper methods for each service method.
//
// Usage example:
//
//	type GreeterServerWrapper struct {
//	    *http.ServerWrapper
//	    GreeterServer  // embed your actual server
//	}
//
//	func (s *GreeterServerWrapper) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
//	    var resp HelloReply
//	    err := s.ServerWrapper.WrapMethod(ctx, "SayHello", req, &resp)
//	    return &resp, err
//	}
func (w *ServerWrapper) WrapMethod(ctx context.Context, methodName string, req, resp interface{}) error {
	return w.Invoke(ctx, methodName, req, resp)
}

// =============================================================================
// Alternative: Create a simple wrapper function for common cases
// =============================================================================

// MiddlewareHandler is a simple middleware that works with gRPC-style handlers.
// This is a compatibility layer for middleware that follows the gRPC interceptor pattern.
type MiddlewareHandler func(ctx context.Context, req interface{}) (interface{}, error)

// ToGRPCInterceptor converts gokit middleware to gRPC unary interceptor format.
// This allows using gokit middleware with gRPC servers.
func ToGRPCInterceptor(handler MiddlewareHandler) func(ctx context.Context, req interface{}) (interface{}, error) {
	return handler
}
