package grpc

import "google.golang.org/grpc"

type Server struct {
	grpcServer *grpc.Server
}

// NewServer 创建 grpc 服务器
func NewServer() *Server {
	return &Server{
		grpcServer: grpc.NewServer(),
	}
}
