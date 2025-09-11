package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/haysons/gokit/transport/testdata/helloworld"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
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
	helloworld.RegisterGreeterServer(server.GetServiceRegistrar(), greeterServer{})

	go func() {
		err := server.Start(ctx)
		require.NoError(t, err)
	}()
	time.Sleep(100 * time.Millisecond)

	req := &helloworld.HelloRequest{
		Name: "hayson",
	}
	resolver.SetDefaultScheme("passthrough")
	conn, err := grpc.NewClient("localhost:8088", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := helloworld.NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "Hello hayson", resp.Message)
}
