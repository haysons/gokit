package errors

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
	"path/filepath"
	"runtime"
	"strings"
)

// WithStack 为错误附加堆栈信息，此堆栈信息会覆盖原本错误中的堆栈信息
// 若错误已经存在堆栈信息，则应当只在需要重新指定错误堆栈时才使用此函数
func WithStack(err error) error {
	return errors.WithStackDepth(err, 1)
}

// GetStackTrace 获取错误中的堆栈信息，若错误为链式错误，则使用最外层错误的堆栈信息，以保证错误堆栈覆盖生效
// errbase.StackTrace 实际上是 pkg/errors中的StackTrace，主流日志库均可友好打印堆栈信息
func GetStackTrace(err error) errbase.StackTrace {
	for err != nil {
		if stp, ok := err.(errbase.StackTraceProvider); ok {
			return stp.StackTrace()
		}
		err = errors.Unwrap(err)
	}
	return nil
}

// GetStack 获取错误中的堆栈信息，以易于日志阅读及解析的字符串表示
func GetStack(err error) []string {
	st := GetStackTrace(err)
	out := make([]string, 0, len(st))
	for _, frame := range st {
		pc := uintptr(frame) - 1
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc)
		shortFile := filepath.Base(file)
		funcName := shortenFuncName(fn.Name())
		out = append(out, fmt.Sprintf("%s:%d %s", shortFile, line, funcName))
	}
	return out
}

func shortenFuncName(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
