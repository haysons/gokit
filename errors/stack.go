package errors

import (
	"fmt"
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
	var allStacks []errbase.StackTrace
	var prev errbase.StackTrace
	for err != nil {
		if stp, ok := err.(errbase.StackTraceProvider); ok {
			st := stp.StackTrace()
			if newSt, elided := errbase.ElideSharedStackTraceSuffix(prev, st); elided {
				st = newSt
			}
			allStacks = append(allStacks, st)
			prev = st
		}
		err = errors.Unwrap(err)
	}

	var finalStack errbase.StackTrace
	for i := len(allStacks) - 1; i >= 0; i-- {
		finalStack = append(finalStack, allStacks[i]...)
	}

	return finalStack
}

// GetSource 获取错误最深处的调用位置，即错误发生的实际位置，若未能获取到调用位置，将返回空
func GetSource(err error) string {
	var innerStack errbase.StackTrace

	for err != nil {
		if stp, ok := err.(errbase.StackTraceProvider); ok {
			innerStack = stp.StackTrace()
		}
		err = errors.Unwrap(err)
	}

	if len(innerStack) > 0 {
		frame := innerStack[0]
		return fmt.Sprintf("%+s:%d", frame, frame)
	}
	return ""
}
