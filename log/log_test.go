package log

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogPrint(t *testing.T) {
	Debug().Msg("debug")
	Info().Msg("info")
	Warn().Msg("warn")
	Error().Msg("error")
	Err(errors.New("error")).Msg("error")
}

func TestFatalFunction(t *testing.T) {
	if os.Getenv("LOG_FATAL_TEST") == "1" {
		Fatal().Msg("fatal from Fatal() function")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatalFunction")
	cmd.Env = append(os.Environ(), "LOG_FATAL_TEST=1")
	output, err := cmd.CombinedOutput()

	// 检查是否是非正常退出
	var exitErr *exec.ExitError
	ok := errors.As(err, &exitErr)
	assert.True(t, ok, "expected process to exit due to Fatal()")
	assert.Equal(t, 1, exitErr.ExitCode(), "expected exit code to be 1 from Fatal()")

	t.Log(string(output))
	assert.Contains(t, string(output), "fatal from Fatal()", "expected fatal message to appear in output")
}

func TestDefaultLogger(t *testing.T) {
	logger := GetDefault()
	assert.Equal(t, zerolog.InfoLevel, logger.GetLevel(), "expected default log level to be info")
}

func TestSetDefault(t *testing.T) {
	conf := &Config{
		Level:        "debug",
		ConsoleFmt:   false,
		ConsoleColor: false,
	}
	SetDefault(conf)
	logger := GetDefault()

	assert.Equal(t, zerolog.DebugLevel, logger.GetLevel(), "expected log level to be debug after SetDefault")
}

func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer

	conf := &Config{
		Level:      "debug",
		ConsoleFmt: false,
	}
	logger := NewZeroLogger(conf).Output(&buf)

	logger.Debug().Msg("debug log")
	logger.Info().Msg("info log")
	logger.Warn().Msg("warn log")
	logger.Error().Msg("error log")

	out := buf.String()
	assert.Contains(t, out, "debug log", "expected debug log message to be present")
	assert.Contains(t, out, "info log", "expected info log message to be present")
	assert.Contains(t, out, "warn log", "expected warn log message to be present")
	assert.Contains(t, out, "error log", "expected error log message to be present")
}

func TestLogErr(t *testing.T) {
	var buf bytes.Buffer
	conf := &Config{
		Level:      "debug",
		ConsoleFmt: false,
	}
	logger := NewZeroLogger(conf).Output(&buf)

	logger.Err(errors.New("test error")).Msg("something happened")

	out := buf.String()
	assert.Contains(t, out, "test error", "expected error message to be present in output")
	assert.Contains(t, out, "something happened", "expected log message to be present in output")
}
