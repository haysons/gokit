package errors

import (
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
)

// WithStack 为错误附加堆栈信息
func WithStack(err error) error {
	return errors.WithStackDepth(err, 1)
}

// GetStackTrace 获取错误中的堆栈信息
// errbase.StackTrace 实际上是 pkg/errors中的StackTrace，主流日志库均可友好打印堆栈信息
func GetStackTrace(err error) errbase.StackTrace {
	var stack errbase.StackTrace
	var prev errbase.StackTrace
	for err != nil {
		if stp, ok := err.(errbase.StackTraceProvider); ok {
			st := stp.StackTrace()
			if newSt, elided := errbase.ElideSharedStackTraceSuffix(prev, st); elided {
				st = newSt
			}
			for i := len(st) - 1; i >= 0; i-- {
				stack = append(stack, st[i])
			}
			prev = st
		}
		err = errors.Unwrap(err)
	}
	return stack
}
