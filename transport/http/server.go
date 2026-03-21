package http

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	gruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/haysons/gokit/middleware"
	"github.com/haysons/gokit/transport"
)

// ServerConfig http 服务器配置项
type ServerConfig struct {
	Addr              string        `mapstructure:"addr"`                // 服务监听的地址，host:port
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`        // 读取请求的超时时间（含header+body）
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"` // 读取请求头的超时时间
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`       // 写响应的超时时间
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`        // Keep-Alive 空闲超时时间

	tlsConf    *tls.Config                    // tls 配置
	muxOptions []gruntime.ServeMuxOption       // grpc-gateway mux 需要用到配置项
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

// WithReadTimeout 配置读取请求（含header+body）的超时时间
func WithReadTimeout(timeout time.Duration) ServerOption {
	return func(c *ServerConfig) {
		c.ReadTimeout = timeout
	}
}

// WithReadHeaderTimeout 配置读取请求头的超时时间
func WithReadHeaderTimeout(timeout time.Duration) ServerOption {
	return func(c *ServerConfig) {
		c.ReadHeaderTimeout = timeout
	}
}

// WithWriteTimeout 配置写响应的超时时间
func WithWriteTimeout(timeout time.Duration) ServerOption {
	return func(c *ServerConfig) {
		c.WriteTimeout = timeout
	}
}

// WithIdleTimeout 配置 Keep-Alive 空闲超时时间
func WithIdleTimeout(timeout time.Duration) ServerOption {
	return func(c *ServerConfig) {
		c.IdleTimeout = timeout
	}
}

// WithTLSConfig 配置 tls 加密相关
func WithTLSConfig(cfg *tls.Config) ServerOption {
	return func(c *ServerConfig) {
		c.tlsConf = cfg
	}
}

// WithMuxOptions 配置 grpc-gateway mux 相关
func WithMuxOptions(opts ...gruntime.ServeMuxOption) ServerOption {
	return func(c *ServerConfig) {
		c.muxOptions = append(c.muxOptions, opts...)
	}
}

// Server http 服务器
//
// 使用方式:
//
//	// 方式 1：使用 protoc 生成的 gateway 代码
//	srv := http.NewServer(http.WithAddr(":8080"))
//	srv.Use(middleware.Logging(), middleware.Recovery())
//	RegisterGreeterHandlerServer(ctx, srv.GetMux(), &myServer{})
//	srv.Start(ctx)
//
//	// 方式 2：手动注册路由（不需要 protoc）
//	srv := http.NewServer(http.WithAddr(":8080"))
//	srv.Use(middleware.Logging(), middleware.Recovery())
//	srv.Handle("/api/hello", SayHello, &HelloRequest{}, &HelloReply{})
//	srv.Start(ctx)
type Server struct {
	httpSrv     *http.Server
	mux         *gruntime.ServeMux
	cfg         *ServerConfig
	middlewares []middleware.Middleware
}

// NewServer 创建 http 服务器
func NewServer(opts ...ServerOption) *Server {
	// 服务配置
	cfg := new(ServerConfig)
	for _, opt := range opts {
		opt(cfg)
	}

	// grpc-gateway mux 配置
	var muxOpts []gruntime.ServeMuxOption
	if len(cfg.muxOptions) > 0 {
		muxOpts = append(muxOpts, cfg.muxOptions...)
	}

	// 构造 http.Server
	mux := gruntime.NewServeMux(muxOpts...)
	httpSrv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		TLSConfig:         cfg.tlsConf,
	}
	return &Server{
		httpSrv:     httpSrv,
		mux:         mux,
		cfg:         cfg,
		middlewares: []middleware.Middleware{},
	}
}

// Use 添加中间件
// 注册的中间件会自动应用到所有 protoc 生成的 gateway 路由和手动注册的路由
//
//	// 示例
//	srv := http.NewServer(http.WithAddr(":8080"))
//	srv.Use(middleware.Logging(), middleware.Recovery())
//	RegisterGreeterHandlerServer(ctx, srv.GetMux(), &myServer{})
func (s *Server) Use(m ...middleware.Middleware) {
	s.middlewares = append(s.middlewares, m...)
	
	// 注册到 grpc-gateway 的全局拦截器
	s.registerGatewayInterceptor()
}

// registerGatewayInterceptor 将 gokit middleware 注册到 grpc-gateway 的拦截器
func (s *Server) registerGatewayInterceptor() {
	if len(s.middlewares) == 0 {
		return
	}
	
	// 创建 gokit middleware 到 grpc-gateway interceptor 的转换
	interceptor := gruntime.Interceptor(func(ctx context.Context, methodName string, req, resp interface{}, handler func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error) {
		// 构建 handler chain
		h := handler
		
		// 应用所有中间件（从外到内）
		for i := len(s.middlewares) - 1; i >= 0; i-- {
			next := h
			m := s.middlewares[i]
			h = func(ctx context.Context, req any) (any, error) {
				return m(func(ctx context.Context, req any) (any, error) {
					return next(ctx, req)
				})(ctx, req)
			}
		}
		
		return h(ctx, req)
	})
	
	// 注册拦截器
	gruntime.RegisterInterceptor(interceptor)
}

// GetMux 获取 ServeMux，主要用于外部注册处理函数
func (s *Server) GetMux() *gruntime.ServeMux {
	return s.mux
}

// Start 启动 http 服务器
func (s *Server) Start(_ context.Context) error {
	var err error
	if s.cfg.tlsConf != nil {
		err = s.httpSrv.ListenAndServeTLS("", "")
	} else {
		err = s.httpSrv.ListenAndServe()
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop 停止 http 服务器
func (s *Server) Stop(ctx context.Context) error {
	// 优先使用 Shutdown 优雅关闭，此方法会判断 ctx
	err := s.httpSrv.Shutdown(ctx)
	if err != nil {
		// 优雅关闭超时，强行关闭
		if ctx.Err() != nil {
			err = s.httpSrv.Close()
		}
	}
	return err
}

// GetTransport 获取当前请求的 transport 信息
// 用于在中间件中获取请求信息
func GetTransport(ctx context.Context) (transport.Transporter, bool) {
	return transport.FromServerContext(ctx)
}
