package log

import (
	"github.com/lmittmann/tint"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

// NewSlogger 创建slog对象
func NewSlogger(conf *Config) *slog.Logger {
	var writer io.Writer

	if conf.Filename == "" {
		writer = os.Stdout
	} else {
		if conf.MaxAge <= 0 {
			conf.MaxAge = 30
		}
		writer = &lumberjack.Logger{
			Filename:  conf.Filename,
			MaxAge:    conf.MaxAge,
			LocalTime: true,
		}
	}

	level := parseLevel(conf.Level)

	var handler slog.Handler
	if conf.ConsoleFmt {
		handler = tint.NewHandler(writer, &tint.Options{
			AddSource:  true,
			Level:      level,
			TimeFormat: time.DateTime,
			NoColor:    !conf.ConsoleColor,
		})
	} else {
		opts := &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}
		handler = slog.NewJSONHandler(writer, opts)
	}

	return slog.New(handler)
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
