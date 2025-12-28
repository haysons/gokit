package http

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/haysons/gokit/middleware"
)

// ServerConfig http 服务器配置项
type ServerConfig struct {
	Addr              string        `mapstructure:"addr"`                // 服务监听的地址，host:port
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`        // 读取请求的超时时间（含header+body）
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"` // 读取请求头的超时时间
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`       // 写响应的超时时间
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`        // Keep-Alive 空闲超时时间

	tlsConf    *tls.Config              // tls 配置
	muxOptions []runtime.ServeMuxOption // grpc-gateway mux 需要用到配置项
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
func WithMuxOptions(opts ...runtime.ServeMuxOption) ServerOption {
	return func(c *ServerConfig) {
		c.muxOptions = opts
	}
}

type Server struct {
	httpSrv    *http.Server
	mux        *runtime.ServeMux
	cfg        *ServerConfig
	middleware []middleware.Middleware
}

// NewServer 创建 http 服务器
func NewServer(opts ...ServerOption) *Server {
	// 服务配置
	cfg := new(ServerConfig)
	for _, opt := range opts {
		opt(cfg)
	}
	srv := &Server{
		cfg: cfg,
	}

	// 中间件处理
	middlewareOpts := runtime.WithMiddlewares(srv.middlewareToMux())

	muxOpts := []runtime.ServeMuxOption{
		middlewareOpts,
	}
	if len(cfg.muxOptions) > 0 {
		muxOpts = append(muxOpts, cfg.muxOptions...)
	}

	// 构造 http.Server
	mux := runtime.NewServeMux(muxOpts...)
	httpSrv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		TLSConfig:         cfg.tlsConf,
	}
	srv.httpSrv = httpSrv
	srv.mux = mux
	return srv
}

// Use 添加中间件
func (s *Server) Use(m ...middleware.Middleware) {
	s.middleware = append(s.middleware, m...)
}

// GetMux 获取 ServeMux，主要用于外部注册处理函数
func (s *Server) GetMux() *runtime.ServeMux {
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

// middlewareToMux 通用中间件转化为 grpc gateway 中间件
func (s *Server) middlewareToMux() runtime.Middleware {
	return func(handlerFunc runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			// todo: 此处考虑能否将 grpc gateway 内部加一段直接使用解码后的 req resp 的中间件的逻辑
			handlerFunc(w, r, pathParams)
		}
	}
}
