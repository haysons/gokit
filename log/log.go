package log

import (
	"log/slog"
)

type Config struct {
	Level        string `mapstructure:"level"`         // 日志级别，支持 debug info warn error 默认info
	Filename     string `mapstructure:"filename"`      // 日志文件名，若文件名为空，则会将日志打印至stdout
	MaxAge       int    `mapstructure:"max_age"`       // 日志文件最大保存天数，默认为30
	ConsoleFmt   bool   `mapstructure:"console_fmt"`   // 日志默认以json方式打印，若ConsoleFmt为true，将以更适合终端阅读的方式打印，此方式性能很差
	ConsoleColor bool   `mapstructure:"console_color"` // 日志采用终端打印格式时是否包含颜色
}

// 默认将日志打印至stdout，使用终端格式打印，且包含颜色，开发体验更好，但性能较差，不适于线上使用
func init() {
	SetDefaultSlog(&Config{
		Level:        "info",
		ConsoleFmt:   true,
		ConsoleColor: true,
	})
}

// GetDefaultSlog 获取默认的日志对象
func GetDefaultSlog() *slog.Logger {
	return slog.Default()
}

// SetDefaultSlog 基于配置信息设置默认的日志对象
func SetDefaultSlog(conf *Config) {
	slog.SetDefault(NewSlogger(conf))
}
