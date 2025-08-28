package config

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/haysons/gokit/log"
	"github.com/spf13/viper"
)

type Config[T any] struct {
	mu     sync.RWMutex
	viper  *viper.Viper
	config T
	logger *slog.Logger
}

// New 新建 Config 实例, T为配置对应的结构体
func New[T any]() *Config[T] {
	return &Config[T]{
		viper:  viper.New(),
		logger: log.GetDefaultSlog(),
	}
}

// SetType 设置配置类型，如：json, yaml, toml
func (c *Config[T]) SetType(t string) {
	c.viper.SetConfigType(t)
}

// SetFile 设置配置文件路径，如：./config.yaml
func (c *Config[T]) SetFile(file string) {
	c.viper.SetConfigFile(file)
}

// AutomaticEnv 自动加载环境变量值作为配置项，如：环境变量A_B_C作为配置项a.b.c的值
func (c *Config[T]) AutomaticEnv() {
	c.viper.AutomaticEnv()
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// SetEnvPrefix 设置环境变量的统一前缀
func (c *Config[T]) SetEnvPrefix(prefix string) {
	c.viper.SetEnvPrefix(prefix)
}

// SetDefault 设置配置项默认值
func (c *Config[T]) SetDefault(key string, value any) {
	c.viper.SetDefault(key, value)
}

// SetLogger 配置日志组件
func (c *Config[T]) SetLogger(logger *slog.Logger) {
	c.logger = logger
}

// Load 加载配置项
func (c *Config[T]) Load() error {
	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}
	return c.unmarshalConfig()
}

func (c *Config[T]) unmarshalConfig() error {
	var cfg T
	if err := c.viper.Unmarshal(&cfg); err != nil {
		return err
	}
	c.mu.Lock()
	c.config = cfg
	c.mu.Unlock()

	c.print()
	return nil
}

// Get 获取当前配置信息
func (c *Config[T]) Get() T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// Watch 监听配置项变化
func (c *Config[T]) Watch() {
	c.viper.WatchConfig()
	c.viper.OnConfigChange(func(e fsnotify.Event) {
		c.logger.Info("config file changed", slog.String("file name", e.Name))
		if err := c.unmarshalConfig(); err != nil {
			c.logger.Error("unmarshal config failed", slog.Any("error", err))
			return
		}
	})
}

func (c *Config[T]) print() {
	for _, k := range c.viper.AllKeys() {
		v := c.viper.Get(k)
		c.logger.Info("config item", slog.String("key", k), slog.Any("value", v))
	}
}
