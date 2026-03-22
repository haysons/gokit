package http

import (
	"context"

	"github.com/haysons/gokit/transport/testdata/helloworld"
)

// =============================================================================
// GreeterServerWrapper - Generated wrapper for Greeter service
// =============================================================================

// GreeterServerWrapper wraps a GreeterServer to add middleware support.
// Use this with protoc-generated RegisterGreeterHandlerServer.
//
// Example:
//
//	// Define your server
//	type MyGreeter struct{}
//	func (s *MyGreeter) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
//	    return &HelloReply{Message: "Hello " + req.Name}, nil
//	}
//
//	// Create wrapped server
//	wrapped := http.NewGreeterServerWrapper(&MyGreeter{})
//	wrapped.Use(func(ctx context.Context, methodName string, req interface{}, handler func(ctx, req) (interface{}, error)) (interface{}, error) {
//	    log.Printf("calling: %s", methodName)
//	    return handler(ctx, req)
//	})
//
//	// Use directly with protoc-generated registration
//	RegisterGreeterHandlerServer(ctx, mux, wrapped)
type GreeterServerWrapper struct {
	helloworld.UnimplementedGreeterServer
	server helloworld.GreeterServer
}

// NewGreeterServerWrapper creates a GreeterServer wrapper with optional interceptors.
func NewGreeterServerWrapper(server helloworld.GreeterServer) *GreeterServerWrapper {
	return &GreeterServerWrapper{
		UnimplementedGreeterServer: helloworld.UnimplementedGreeterServer{},
		server: server,
	}
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
		return w.server.SayHello(ctx, req.(*helloworld.HelloRequest))
	})
	if err != nil {
		return nil, err
	}
	return result.(*helloworld.HelloReply), nil
}
