package http

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/haysons/gokit/transport/testdata/helloworld"
	"github.com/stretchr/testify/require"
)

type greeterServer struct {
	helloworld.UnimplementedGreeterServer
}

func (g greeterServer) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + request.Name}, nil
}

func TestServerStart(t *testing.T) {
	ctx := context.Background()
	server := NewServer(WithAddr(":8088"))
	err := helloworld.RegisterGreeterHandlerServer(ctx, server.GetMux(), greeterServer{})
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
}
