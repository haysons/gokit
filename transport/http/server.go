package http

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Server struct {
	mux *runtime.ServeMux
}

func NewServer() *Server {
	return &Server{
		mux: runtime.NewServeMux(),
	}
}
