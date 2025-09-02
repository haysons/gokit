package http

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// ServerConfig http 服务器配置项
type ServerConfig struct {
	Addr              string        `mapstructure:"addr"`                // 服务监听的地址，host:port
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`        // 读取请求的超时时间（含header+body）
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"` // 读取请求头的超时时间
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`       // 写响应的超时时间
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`        // Keep-Alive 空闲超时时间

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

// WithMuxOptions 配置 grpc-gateway mux 相关
func WithMuxOptions(opts ...runtime.ServeMuxOption) ServerOption {
	return func(c *ServerConfig) {
		c.muxOptions = opts
	}
}

type Server struct {
	httpSrv *http.Server
	mux     *runtime.ServeMux
}

// NewServer 创建 http 服务器
func NewServer(opts ...ServerOption) *Server {
	// 服务配置
	cfg := new(ServerConfig)
	for _, opt := range opts {
		opt(cfg)
	}

	// 构造 http.Server
	mux := runtime.NewServeMux(cfg.muxOptions...)
	httpSrv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	return &Server{
		httpSrv: httpSrv,
		mux:     mux,
	}
}

// GetMux 获取 ServeMux，主要用于外部注册处理函数
func (s *Server) GetMux() *runtime.ServeMux {
	return s.mux
}

// Start 启动 http 服务器
func (s *Server) Start(ctx context.Context) error {
	return s.httpSrv.ListenAndServe()
}

// Stop 停止 http 服务器
func (s *Server) Stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
