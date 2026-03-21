package http

import (
	"context"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/haysons/gokit/middleware"
	"github.com/haysons/gokit/transport/testdata/helloworld"
	"github.com/stretchr/testify/require"
)

type testGreeterServer struct {
	helloworld.UnimplementedGreeterServer
}

func (g testGreeterServer) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + request.Name}, nil
}

// TestServerStart tests basic HTTP server functionality with protoc-generated gateway
//
// Note: This test uses the original protoc-generated gateway code.
// For middleware to work with protoc-generated handlers, you need to:
// 1. Use the modified protoc-gen-grpc-gateway from github.com/haysons/grpc-gateway
// 2. Regenerate your .pb.gw.go files with the modified plugin
func TestServerStart(t *testing.T) {
	ctx := context.Background()
	server := NewServer(WithAddr(":8088"))
	err := helloworld.RegisterGreeterHandlerServer(ctx, server.GetMux(), testGreeterServer{})
	require.NoError(t, err)

	go func() {
		err = server.Start(ctx)
		require.NoError(t, err)
	}()
	time.Sleep(100 * time.Millisecond)

	req := &helloworld.HelloRequest{
		Name: "hayson",
	}
	reply := &helloworld.HelloReply{}
	res, err := resty.New().R().
		SetBody(req).
		SetResult(reply).
		Post("http://localhost:8088/v1/hello")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Equal(t, "Hello hayson", reply.Message)

	// cleanup
	server.Stop(context.Background())
}

// TestServerWithMiddleware tests that middlewares work with manually registered handlers
//
// For middlewares to work with protoc-generated gateway handlers,
// you need to regenerate your .pb.gw.go files using the modified protoc-gen-grpc-gateway
// from github.com/haysons/grpc-gateway
func TestServerWithMiddleware(t *testing.T) {
	// Track middleware execution
	var middlewareCalled atomic.Bool
	var capturedName string

	// Create a custom middleware
	testMiddleware := func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			middlewareCalled.Store(true)
			// Capture request data
			if r, ok := req.(*helloworld.HelloRequest); ok {
				capturedName = r.Name
			}
			return next(ctx, req)
		}
	}

	ctx := context.Background()
	server := NewServer(WithAddr(":8089"))

	// Register middleware
	server.Use(testMiddleware)

	// Register handler manually (not using protoc-generated code)
	// For protoc-generated handlers to work with middleware, use the modified protoc-gen-grpc-gateway
	server.RegisterHandler("POST", "/api/hello", 
		func(ctx context.Context, req interface{}) (interface{}, error) {
			r, ok := req.(*helloworld.HelloRequest)
			if !ok {
				return nil, nil
			}
			return &helloworld.HelloReply{Message: "Hello " + r.Name}, nil
		},
		&helloworld.HelloRequest{},
		&helloworld.HelloReply{},
	)

	go func() {
		err := server.Start(ctx)
		require.NoError(t, err)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make request
	req := &helloworld.HelloRequest{
		Name: "middleware-test",
	}
	reply := &helloworld.HelloReply{}
	res, err := resty.New().R().
		SetBody(req).
		SetResult(reply).
		Post("http://localhost:8089/api/hello")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Equal(t, "Hello middleware-test", reply.Message)

	// Verify middleware was called and captured the request
	require.True(t, middlewareCalled.Load(), "Middleware should have been called")
	require.Equal(t, "middleware-test", capturedName)

	// cleanup
	server.Stop(context.Background())
}

// TestServerWithMultipleMiddlewares tests multiple middleware chain with manual registration
func TestServerWithMultipleMiddlewares(t *testing.T) {
	callOrder := make([]string, 0, 3)

	// Middleware 1
	middleware1 := func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			callOrder = append(callOrder, "1_before")
			resp, err := next(ctx, req)
			callOrder = append(callOrder, "1_after")
			return resp, err
		}
	}

	// Middleware 2
	middleware2 := func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			callOrder = append(callOrder, "2_before")
			resp, err := next(ctx, req)
			callOrder = append(callOrder, "2_after")
			return resp, err
		}
	}

	ctx := context.Background()
	server := NewServer(WithAddr(":8090"))

	// Register middlewares in order
	server.Use(middleware1)
	server.Use(middleware2)

	// Register handler manually
	server.RegisterHandler("POST", "/api/chain", 
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return &helloworld.HelloReply{Message: "OK"}, nil
		},
		&helloworld.HelloRequest{},
		&helloworld.HelloReply{},
	)

	go func() {
		err := server.Start(ctx)
		require.NoError(t, err)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make request
	res, err := resty.New().R().
		SetBody(`{"name": "chain"}`).
		Post("http://localhost:8090/api/chain")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())

	// Verify middleware was called (at least one)
	require.GreaterOrEqual(t, len(callOrder), 2, "At least some middlewares should be called")
	t.Logf("Call order: %v", callOrder)

	// cleanup
	server.Stop(context.Background())
}

// TestServerManualRegistration tests manual handler registration (without protoc)
func TestServerManualRegistration(t *testing.T) {
	ctx := context.Background()
	server := NewServer(WithAddr(":8091"))

	// Track middleware execution
	var middlewareCalled atomic.Bool
	testMiddleware := func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			middlewareCalled.Store(true)
			return next(ctx, req)
		}
	}
	server.Use(testMiddleware)

	// Manually register handler (without protoc)
	server.RegisterHandler("POST", "/api/hello", 
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return &helloworld.HelloReply{Message: "Hello "}, nil
		},
		&helloworld.HelloRequest{},
		&helloworld.HelloReply{},
	)

	go func() {
		err := server.Start(ctx)
		require.NoError(t, err)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make request
	res, err := resty.New().R().
		SetBody(`{"name": "manual"}`).
		Post("http://localhost:8091/api/hello")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.True(t, middlewareCalled.Load(), "Middleware should be called for manually registered handler")
	require.True(t, strings.Contains(string(res.Body()), "Hello"), "Response should contain message")

	// cleanup
	server.Stop(context.Background())
}
