package grpc

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/haysons/gokit/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// ServerConfig grpc 服务配置项
type ServerConfig struct {
	Addr string `mapstructure:"addr"` // 服务监听的地址，host:port

	tlsConf           *tls.Config                    // tls 配置
	unaryInts         []grpc.UnaryServerInterceptor  // grpc 请求响应式拦截器
	streamInts        []grpc.StreamServerInterceptor // grpc 流式拦截器
	grpcOpts          []grpc.ServerOption            // grpc 原生配置
	disableReflection bool                           // 关闭服务端反射（服务端反射可以为客户端提供服务器有哪些服务及方法）
}

// ServerOption 函数式配置项
type ServerOption func(*ServerConfig)

// WithConfig 整体替换配置
func WithConfig(cfg ServerConfig) ServerOption {
	return func(c *ServerConfig) {
		*c = cfg
	}
}

// WithAddr 配置服务监听地址
func WithAddr(addr string) ServerOption {
	return func(c *ServerConfig) {
		c.Addr = addr
	}
}

// WithTLSConfig 配置 tls 加密相关
func WithTLSConfig(cfg *tls.Config) ServerOption {
	return func(c *ServerConfig) {
		c.tlsConf = cfg
	}
}

// WithUnaryInterceptor 配置 grpc 请求响应式拦截器
func WithUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(c *ServerConfig) {
		c.unaryInts = in
	}
}

// WithStreamInterceptor 配置 grpc 流式拦截器
func WithStreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(c *ServerConfig) {
		c.streamInts = in
	}
}

// WithDisableReflection 禁用服务端反射
func WithDisableReflection() ServerOption {
	return func(c *ServerConfig) {
		c.disableReflection = true
	}
}

// WithGRPCOptions 增加原生 grpc 配置
func WithGRPCOptions(opts ...grpc.ServerOption) ServerOption {
	return func(c *ServerConfig) {
		c.grpcOpts = opts
	}
}

type Server struct {
	grpcServer       *grpc.Server
	cfg              *ServerConfig
	middleware       []middleware.Middleware
	streamMiddleware []middleware.Middleware
}

// NewServer 创建 grpc 服务器
func NewServer(opts ...ServerOption) *Server {
	// 服务配置
	cfg := new(ServerConfig)
	for _, opt := range opts {
		opt(cfg)
	}
	srv := &Server{
		cfg: cfg,
	}

	// 拦截器
	unaryInts := []grpc.UnaryServerInterceptor{
		srv.middlewareToUnaryInterceptor(),
	}
	streamInts := make([]grpc.StreamServerInterceptor, 0, len(cfg.streamInts))
	if len(cfg.unaryInts) > 0 {
		unaryInts = append(unaryInts, cfg.unaryInts...)
	}
	if len(cfg.streamInts) > 0 {
		streamInts = append(streamInts, cfg.streamInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}

	// tls 配置
	if cfg.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(cfg.tlsConf)))
	}

	// grpc 原生配置
	if len(cfg.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, cfg.grpcOpts...)
	}

	grpcServer := grpc.NewServer(grpcOpts...)

	// 是否禁用服务端反射
	if !cfg.disableReflection {
		reflection.Register(grpcServer)
	}

	srv.grpcServer = grpcServer
	return srv
}

// Use 请求响应式接口添加中间件
func (s *Server) Use(m ...middleware.Middleware) {
	s.middleware = append(s.middleware, m...)
}

// UseStream 流式接口添加中间件
func (s *Server) UseStream(m ...middleware.Middleware) {
	s.streamMiddleware = append(s.streamMiddleware, m...)
}

// GetServiceRegistrar 注册 grpc service
func (s *Server) GetServiceRegistrar() grpc.ServiceRegistrar {
	return s.grpcServer
}

// Start 启动 grpc 服务器
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(lis)
}

// Stop 停止 grpc 服务器
func (s *Server) Stop(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		// 优雅关闭
		defer close(done)
		s.grpcServer.GracefulStop()
	}()

	select {
	case <-done:
	case <-ctx.Done():
		// 强行关闭
		s.grpcServer.Stop()
	}
	return nil
}

// middlewareToUnaryInterceptor 通用中间件转化为 grpc 拦截器
func (s *Server) middlewareToUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		h := func(ctx context.Context, req any) (any, error) {
			return handler(ctx, req)
		}
		h = middleware.Combine(s.middleware...)(h)
		return h(ctx, req)
	}
}
