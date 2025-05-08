package errors

import (
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
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
