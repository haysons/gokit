package log

import (
	"github.com/rs/zerolog"
)

type Config struct {
	Level        string `yaml:"level"`         // 日志级别，支持 debug info warn error 默认info
	Filename     string `yaml:"filename"`      // 日志文件名，若文件名为空，则会将日志打印至stdout
	MaxAge       int    `yaml:"max_age"`       // 日志文件最大保存天数，默认为30
	ConsoleFmt   bool   `yaml:"console_fmt"`   // 日志默认以json方式打印，若ConsoleFmt为true，将以更适合终端阅读的方式打印，此方式性能很差
	ConsoleColor bool   `yaml:"console_color"` // 日志采用终端打印格式时是否包含颜色
}

// defaultLogger 默认将日志打印至stdout，使用终端格式打印，且包含颜色，开发体验更好，但性能较差，不适于线上使用
var defaultLogger = NewZeroLogger(&Config{
	Level:        "info",
	ConsoleFmt:   true,
	ConsoleColor: true,
})

// GetDefault 获取默认的日志对象
func GetDefault() zerolog.Logger {
	return defaultLogger
}

// SetDefault 基于配置信息设置默认的日志对象
func SetDefault(conf *Config) {
	defaultLogger = NewZeroLogger(conf)
}

// Debug 打印debug级别日志
func Debug() *zerolog.Event {
	return defaultLogger.Debug()
}

// Info 打印info级别日志
func Info() *zerolog.Event {
	return defaultLogger.Info()
}

// Warn 打印warn级别日志
func Warn() *zerolog.Event {
	return defaultLogger.Warn()
}

// Error 打印error级别日志
func Error() *zerolog.Event {
	return defaultLogger.Error()
}

// Err 基于err快捷打印一个error日志
func Err(err error) *zerolog.Event {
	return defaultLogger.Err(err)
}

// Fatal 打印fatal级别日志
func Fatal() *zerolog.Event {
	return defaultLogger.Fatal()
}
