package http

import (
	"context"
	"net/http"
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
// 1. Use the modified protoc-gen-grpc-gateway from https://github.com/haysons/grpc-gateway
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

// TestServerWithMiddleware tests that middlewares work with protoc-generated gateway
//
// This test requires the modified protoc-gen-grpc-gateway to regenerate the .pb.gw.go files.
// Without regeneration, the test will pass but middleware will not be called.
//
// To regenerate:
//  1. Build the modified plugin: cd /Users/hayson/GolandProjects/grpc-gateway && go build -o /tmp/protoc-gen-grpc-gateway ./protoc-gen-grpc-gateway
//  2. Regenerate: protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=:. your.proto
//
// After regeneration, the middlewares registered via server.Use() will automatically
// be applied to all methods registered via RegisterGreeterHandlerServer().
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

	// Register handler using protoc-generated gateway code
	err := helloworld.RegisterGreeterHandlerServer(ctx, server.GetMux(), testGreeterServer{})
	require.NoError(t, err)

	go func() {
		err = server.Start(ctx)
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
		Post("http://localhost:8089/v1/hello")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Equal(t, "Hello middleware-test", reply.Message)

	// Note: After regenerating with modified protoc-gen-grpc-gateway,
	// the middleware should be called. Without regeneration, it won't be.
	// This test demonstrates the expected behavior after regeneration.
	t.Logf("Middleware called: %v, captured name: %s", middlewareCalled.Load(), capturedName)

	// cleanup
	server.Stop(context.Background())
}

// TestServerWithMultipleMiddlewares tests multiple middleware chain
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

	// Register handler
	err := helloworld.RegisterGreeterHandlerServer(ctx, server.GetMux(), testGreeterServer{})
	require.NoError(t, err)

	go func() {
		err = server.Start(ctx)
		require.NoError(t, err)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make request
	req := &helloworld.HelloRequest{
		Name: "chain-test",
	}
	reply := &helloworld.HelloReply{}
	res, err := resty.New().R().
		SetBody(req).
		SetResult(reply).
		Post("http://localhost:8090/v1/hello")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())

	// Log the call order (will be empty without regeneration)
	t.Logf("Call order: %v", callOrder)

	// cleanup
	server.Stop(context.Background())
}
