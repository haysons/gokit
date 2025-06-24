package log

import (
	"context"
	"log/slog"
)

type Config struct {
	Level        string `yaml:"level"`         // 日志级别，支持 debug info warn error 默认info
	Filename     string `yaml:"filename"`      // 日志文件名，若文件名为空，则会将日志打印至stdout
	MaxAge       int    `yaml:"max_age"`       // 日志文件最大保存天数，默认为30
	ConsoleFmt   bool   `yaml:"console_fmt"`   // 日志默认以json方式打印，若ConsoleFmt为true，将以更适合终端阅读的方式打印，此方式性能很差
	ConsoleColor bool   `yaml:"console_color"` // 日志采用终端打印格式时是否包含颜色
}

// 默认将日志打印至stdout，使用终端格式打印，且包含颜色，开发体验更好，但性能较差，不适于线上使用
func init() {
	SetDefault(&Config{
		Level:        "info",
		ConsoleFmt:   true,
		ConsoleColor: true,
	})
}

// GetDefault 获取默认的日志对象
func GetDefault() *slog.Logger {
	return slog.Default()
}

// SetDefault 基于配置信息设置默认的日志对象
func SetDefault(conf *Config) {
	slog.SetDefault(NewSlogger(conf))
}

// With 为日志附加通用属性
func With(args ...any) *slog.Logger {
	return slog.With(args...)
}

// Debug 打印debug日志
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// DebugCtx 打印debug日志
func DebugCtx(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

// Info 打印info日志
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// InfoCtx 打印info日志
func InfoCtx(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

// Warn 打印warn日志
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// WarnCtx 打印warn日志
func WarnCtx(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

// Error 打印error日志
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// ErrorCtx 打印error日志
func ErrorCtx(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}
