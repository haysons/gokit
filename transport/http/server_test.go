package http

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
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
func TestServerWithMiddleware(t *testing.T) {
	// Clear any previous interceptors
	helloworld.ClearGokitInterceptors()
	defer helloworld.ClearGokitInterceptors()

	// Track middleware execution
	var middlewareCalled atomic.Bool
	var capturedName string

	// Create a custom middleware using gokit interceptor
	testMiddleware := helloworld.GokitInterceptor(func(ctx context.Context, methodName string, req interface{}, handler func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error) {
		middlewareCalled.Store(true)
		// Capture request data
		if r, ok := req.(*helloworld.HelloRequest); ok {
			capturedName = r.Name
		}
		return handler(ctx, req)
	})

	// Register middleware
	helloworld.RegisterGokitInterceptor(testMiddleware)

	ctx := context.Background()
	server := NewServer(WithAddr(":8089"))
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

	// Verify middleware was called and captured the request
	require.True(t, middlewareCalled.Load(), "Middleware should have been called")
	require.Equal(t, "middleware-test", capturedName)

	// cleanup
	server.Stop(context.Background())
}

// TestServerWithMultipleMiddlewares tests multiple middleware chain
func TestServerWithMultipleMiddlewares(t *testing.T) {
	// Clear any previous interceptors
	helloworld.ClearGokitInterceptors()
	defer helloworld.ClearGokitInterceptors()

	callOrder := make([]string, 0, 3)

	// Middleware 1
	middleware1 := helloworld.GokitInterceptor(func(ctx context.Context, methodName string, req interface{}, handler func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error) {
		callOrder = append(callOrder, "1_before")
		r, err := handler(ctx, req)
		callOrder = append(callOrder, "1_after")
		return r, err
	})

	// Middleware 2
	middleware2 := helloworld.GokitInterceptor(func(ctx context.Context, methodName string, req interface{}, handler func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error) {
		callOrder = append(callOrder, "2_before")
		r, err := handler(ctx, req)
		callOrder = append(callOrder, "2_after")
		return r, err
	})

	// Register middlewares
	helloworld.RegisterGokitInterceptor(middleware1)
	helloworld.RegisterGokitInterceptor(middleware2)

	ctx := context.Background()
	server := NewServer(WithAddr(":8090"))
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

	// Verify middleware execution order
	// Middlewares are executed in reverse registration order: last registered = outermost (executed first)
	expectedOrder := []string{
		"2_before",   // middleware2 is registered last, so it's outermost
		"1_before",   // middleware1 is inner
		"1_after",    // then returns
		"2_after",    // then outermost returns
	}
	require.Equal(t, expectedOrder, callOrder, "Middleware chain should execute in order")
	t.Logf("Call order: %v", callOrder)

	// cleanup
	server.Stop(context.Background())
}
