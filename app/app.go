package app

import (
	"context"

	"go.uber.org/fx"
)

type Config struct {
	Name    string `mapstructure:"name"`    // app 名称
	Version string `mapstructure:"version"` // app 版本
	Commit  string `mapstructure:"commit"`  // 程序提交号

	fxOptions []fx.Option // uber fx 配置项
}

// Option 函数式配置项
type Option func(*Config)

// WithConfig 整体替换配置
func WithConfig(cfg Config) Option {
	return func(c *Config) {
		*c = cfg
	}
}

// WithName 配置应用名称
func WithName(n string) Option {
	return func(c *Config) {
		c.Name = n
	}
}

// WithVersion 配置应用版本号
func WithVersion(v string) Option {
	return func(c *Config) {
		c.Version = v
	}
}

// WithCommit 配置应用提交号
func WithCommit(commit string) Option {
	return func(c *Config) {
		c.Commit = commit
	}
}

// WithProvides 配置服务提供者，用于依赖注入
func WithProvides(provide ...fx.Option) Option {
	return func(c *Config) {
		c.fxOptions = append(c.fxOptions)
	}
}

// App 管理整个应用程序的生命周期，基于 uber fx 实现，以此解决依赖注入问题
type App struct {
	cfg   *Config
	fxApp *fx.App
}

// New 创建一个新的 App 实例
func New(opts ...Option) *App {
	// 应用配置项
	cfg := new(Config)
	for _, opt := range opts {
		opt(cfg)
	}

	return &App{
		fxApp: fx.New(fx.Options(cfg.fxOptions...)),
	}
}

// Run 启动应用程序
func (a *App) Run() {
	a.fxApp.Run()
}

// Start 启动应用程序但不阻塞
func (a *App) Start(ctx context.Context) error {
	return a.fxApp.Start(ctx)
}

// Stop 停止应用程序
func (a *App) Stop(ctx context.Context) error {
	return a.fxApp.Stop(ctx)
}
