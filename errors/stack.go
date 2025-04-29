package errors

import (
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
)

// WithStack 为错误附加堆栈信息
func WithStack(err error) error {
	return errors.WithStack(err)
}

// GetStackTrace 获取错误中的堆栈信息
// errbase.StackTrace 实际上是 pkg/errors中的StackTrace，主流日志库均可友好打印堆栈信息
func GetStackTrace(err error) errbase.StackTrace {
	if st, ok := err.(errbase.StackTraceProvider); ok {
		return st.StackTrace()
	}
	return nil
}
