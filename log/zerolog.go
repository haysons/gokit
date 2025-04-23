package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func NewZeroLogger(conf *Config) zerolog.Logger {
	// 日志输出位置
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
	if conf.ConsoleFmt {
		// ConsoleWriter性能很差，console writer并不是encoder而是writer，实际上日志信息都是以json的方式进行序列化的，
		// ConsoleWriter获得bytes之后，需要对其再进行json反序列化，然后再按console writer的方式输出日志，故性能很差
		writer = zerolog.ConsoleWriter{Out: writer, NoColor: !conf.ConsoleColor, TimeFormat: time.DateTime}
	}

	// 日志级别
	level, err := zerolog.ParseLevel(conf.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	// 全局配置项
	zerolog.TimeFieldFormat = time.DateTime
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		dirs := strings.Split(file, string(filepath.Separator))
		n := len(dirs)
		if n >= 2 {
			file = filepath.Join(dirs[n-2], dirs[n-1])
		}
		return file + ":" + strconv.Itoa(line)
	}
	logger := zerolog.New(writer).Level(level).
		Hook(ctxHook{}).With(). // 打印ctx中元数据
		Timestamp().            // 打印日志时间
		Stack().                // 发生错误时打印堆栈信息
		Caller().               // 打印调用位置
		Logger()
	return logger
}

type ctxHook struct{}

func (ch ctxHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	ctx := e.GetCtx()
	reqID := ctx.Value("request_id")
	if reqID != nil {
		reqIDStr, ok := reqID.(string)
		if ok {
			e.Str("request_id", reqIDStr)
		}
	}
}
