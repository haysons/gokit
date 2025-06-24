package log

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestNewSlogger_JSON(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	conf := &Config{
		Filename:     "",
		Level:        "info",
		ConsoleFmt:   false,
		ConsoleColor: false,
	}
	logger := NewSlogger(conf)
	logger.InfoContext(context.Background(), "json log", slog.String("key", "value"))

	w.Close()
	var outBuf bytes.Buffer
	_, _ = outBuf.ReadFrom(r)
	output := outBuf.String()

	assert.Contains(t, output, `"msg":"json log"`, "should contain message")
	assert.Contains(t, output, `"key":"value"`, "should contain key attribute")
}

func TestNewSlogger_ConsoleFmt(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	conf := &Config{
		Filename:     "",
		Level:        "debug",
		ConsoleFmt:   true,
		ConsoleColor: false,
	}
	logger := NewSlogger(conf)
	logger.Debug("debug log", slog.String("foo", "bar"))

	w.Close()
	var outBuf bytes.Buffer
	_, _ = outBuf.ReadFrom(r)
	output := outBuf.String()

	t.Log(output)
	assert.Contains(t, output, "DBG", "should contain DEBUG level")
	assert.Contains(t, output, "debug log", "should contain debug message")
	assert.Contains(t, output, "foo=bar", "should contain custom field")
}

func TestSetDefaultAndGetDefault(t *testing.T) {
	conf := &Config{
		Level:        "info",
		ConsoleFmt:   true,
		ConsoleColor: true,
	}
	SetDefault(conf)
	ctx := context.Background()
	slog.Debug("debug log", slog.String("foo", "bar"))
	slog.Info("info log", slog.String("foo", "bar"))
	slog.Warn("warn log", slog.String("foo", "bar"))
	slog.Error("error log", slog.String("foo", "bar"))
	Debug("debug log", slog.String("foo", "bar"))
	DebugCtx(ctx, "debug log", slog.String("foo", "bar"))
	Info("info log", slog.String("foo", "bar"))
	InfoCtx(ctx, "info log", slog.String("foo", "bar"))
	Warn("warn log", slog.String("foo", "bar"))
	WarnCtx(ctx, "warn log", slog.String("foo", "bar"))
	Error("error log", slog.String("foo", "bar"))
	ErrorCtx(ctx, "error log", slog.String("foo", "bar"))

	logger := With(slog.String("common", "bar"))
	logger.Info("info log", slog.String("foo", "bar"))
}

func TestParseLevel(t *testing.T) {
	assert.Equal(t, slog.LevelDebug, parseLevel("debug"))
	assert.Equal(t, slog.LevelInfo, parseLevel("info"))
	assert.Equal(t, slog.LevelWarn, parseLevel("warn"))
	assert.Equal(t, slog.LevelError, parseLevel("error"))
	assert.Equal(t, slog.LevelInfo, parseLevel("invalid")) // 默认值
}
