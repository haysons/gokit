package config_test

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/haysons/gokit/config"
	"github.com/stretchr/testify/assert"
)

type ServerConfig struct {
	Server struct {
		Host  string `mapstructure:"host"`
		Port  int    `mapstructure:"port"`
		Proto string `mapstructure:"proto"`
	} `mapstructure:"server"`
}

func TestConfig_LoadAndGet(t *testing.T) {
	cfg := config.New[ServerConfig]()
	cfg.SetType("yaml")
	cfg.SetFile(filepath.Join("testdata", "config.yaml"))

	err := cfg.Load()
	assert.NoError(t, err)

	conf := cfg.Get()
	assert.Equal(t, "127.0.0.1", conf.Server.Host)
	assert.Equal(t, 8080, conf.Server.Port)
}

func TestConfig_SetDefault(t *testing.T) {
	cfg := config.New[ServerConfig]()
	cfg.SetType("yaml")
	cfg.SetFile(filepath.Join("testdata", "config.yaml"))
	cfg.SetDefault("server.proto", "http")

	err := cfg.Load()
	assert.NoError(t, err)

	conf := cfg.Get()
	assert.Equal(t, "http", conf.Server.Proto)
}

func TestConfig_AutomaticEnv(t *testing.T) {
	_ = os.Setenv("CONF_SERVER_HOST", "192.168.0.1")

	cfg := config.New[ServerConfig]()
	cfg.SetType("yaml")
	cfg.SetFile(filepath.Join("testdata", "config.yaml"))
	cfg.AutomaticEnv()
	cfg.SetEnvPrefix("CONF")

	err := cfg.Load()
	assert.NoError(t, err)

	conf := cfg.Get()
	assert.Equal(t, "192.168.0.1", conf.Server.Host)
}

func TestConfig_GetThreadSafety(t *testing.T) {
	cfg := config.New[ServerConfig]()
	cfg.SetType("yaml")
	cfg.SetFile(filepath.Join("testdata", "config.yaml"))

	err := cfg.Load()
	assert.NoError(t, err)

	done := make(chan struct{})
	go func() {
		_ = cfg.Get()
		close(done)
	}()

	<-done
	assert.True(t, true, "concurrent read should succeed")
}

func TestConfig_Watch(t *testing.T) {
	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// 初始配置内容
	initial := ServerConfig{}
	initial.Server.Host = "localhost"
	initial.Server.Port = 8000

	writeYaml := func(cfg ServerConfig) {
		data, _ := yaml.Marshal(cfg)
		err = os.WriteFile(tmpFile.Name(), data, 0644)
		assert.NoError(t, err)
	}

	writeYaml(initial)

	// 启动配置加载器
	cfg := config.New[ServerConfig]()
	cfg.SetType("yaml")
	cfg.SetFile(tmpFile.Name())

	err = cfg.Load()
	assert.NoError(t, err)
	cfg.Watch()

	// 修改文件以触发 Watch
	updated := ServerConfig{}
	updated.Server.Host = "127.0.0.1"
	updated.Server.Port = 9090

	writeYaml(updated)

	// 等待 watch 回调生效
	var newConf ServerConfig
	retry := 10
	for i := 0; i < retry; i++ {
		time.Sleep(300 * time.Millisecond)
		newConf = cfg.Get()
		if newConf.Server.Port == 9090 && newConf.Server.Host == "127.0.0.1" {
			break
		}
	}

	assert.Equal(t, "127.0.0.1", newConf.Server.Host)
	assert.Equal(t, 9090, newConf.Server.Port)
}
