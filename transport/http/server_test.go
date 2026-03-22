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

// TestServerWithMiddleware tests that middlewares work with GreeterServerWrapper
func TestServerWithMiddleware(t *testing.T) {
	// Clear previous interceptors
	ClearInterceptors()
	defer ClearInterceptors()

	// Track middleware execution
	var middlewareCalled atomic.Bool
	var capturedName string

	// Create an interceptor
	testInterceptor := Interceptor(func(ctx context.Context, methodName string, req interface{}, handler func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error) {
		middlewareCalled.Store(true)
		// Capture request data
		if r, ok := req.(*helloworld.HelloRequest); ok {
			capturedName = r.Name
		}
		return handler(ctx, req)
	})

	// Create wrapped server with interceptors
	wrapped := NewGreeterServerWrapper(&testGreeterServer{})
	wrapped.Use(testInterceptor)

	ctx := context.Background()
	server := NewServer(WithAddr(":8089"))
	err := helloworld.RegisterGreeterHandlerServer(ctx, server.GetMux(), wrapped)
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
	// Clear previous interceptors
	ClearInterceptors()
	defer ClearInterceptors()

	callOrder := make([]string, 0, 3)

	// Interceptor 1
	interceptor1 := Interceptor(func(ctx context.Context, methodName string, req interface{}, handler func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error) {
		callOrder = append(callOrder, "1_before")
		r, err := handler(ctx, req)
		callOrder = append(callOrder, "1_after")
		return r, err
	})

	// Interceptor 2
	interceptor2 := Interceptor(func(ctx context.Context, methodName string, req interface{}, handler func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error) {
		callOrder = append(callOrder, "2_before")
		r, err := handler(ctx, req)
		callOrder = append(callOrder, "2_after")
		return r, err
	})

	// Create wrapped server with interceptors (last registered = outermost)
	wrapped := NewGreeterServerWrapper(&testGreeterServer{})
	wrapped.Use(interceptor1)
	wrapped.Use(interceptor2)

	ctx := context.Background()
	server := NewServer(WithAddr(":8090"))
	err := helloworld.RegisterGreeterHandlerServer(ctx, server.GetMux(), wrapped)
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

	// Verify interceptor execution order (first registered = outermost)
	expectedOrder := []string{
		"1_before", // interceptor1 is registered first, so it's outermost
		"2_before", // interceptor2 is inner
		"2_after",  // then returns
		"1_after",  // then outermost returns
	}
	require.Equal(t, expectedOrder, callOrder, "Interceptor chain should execute in order")
	t.Logf("Call order: %v", callOrder)

	// cleanup
	server.Stop(context.Background())
}

// TestServerWithoutMiddleware tests that server works without middleware
func TestServerWithoutMiddleware(t *testing.T) {
	// Clear previous interceptors
	ClearInterceptors()
	defer ClearInterceptors()

	// Create wrapped server without interceptors
	wrapped := NewGreeterServerWrapper(&testGreeterServer{})

	ctx := context.Background()
	server := NewServer(WithAddr(":8091"))
	err := helloworld.RegisterGreeterHandlerServer(ctx, server.GetMux(), wrapped)
	require.NoError(t, err)

	go func() {
		err = server.Start(ctx)
		require.NoError(t, err)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make request
	req := &helloworld.HelloRequest{Name: "no-middleware"}
	reply := &helloworld.HelloReply{}
	res, err := resty.New().R().
		SetBody(req).
		SetResult(reply).
		Post("http://localhost:8091/v1/hello")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Equal(t, "Hello no-middleware", reply.Message)

	// cleanup
	server.Stop(context.Background())
}
