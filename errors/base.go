package errors

import "github.com/cockroachdb/errors"

// New 创建一个错误，此错误携带堆栈信息
func New(msg string) error {
	return errors.NewWithDepth(1, msg)
}

// Newf 使用格式化方式创建一个错误，此错误携带堆栈信息
func Newf(format string, args ...any) error {
	return errors.NewWithDepthf(1, format, args...)
}

// Is 判断错误是否为相同错误，当错误为链式错误时，错误链中任意错误和目标错误相等则视作相等
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As 将错误转换为特定类型错误，当错误为链式错误时，错误链中任意错误可转换为目标错误则完成转换
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Wrap 使用msg作为前缀包装错误，并附加堆栈信息，此堆栈信息会覆盖原本错误中的堆栈信息
// 若仅是为了给错误添加上下文信息，应当使用 WithMessage
func Wrap(err error, msg string) error {
	return errors.WrapWithDepth(1, err, msg)
}

// Wrapf 使用格式化字符串作为前缀包装错误，并附加堆栈信息，此堆栈信息会覆盖原本错误中的堆栈信息
// 若仅是为了给错误添加上下文信息，应当使用 WithMessagef
func Wrapf(err error, format string, args ...any) error {
	return errors.WrapWithDepthf(1, err, format, args...)
}

// WithMessage 为错误附加上下文信息
func WithMessage(err error, msg string) error {
	return errors.WithMessage(err, msg)
}

// WithMessagef 使用格式化字符串为错误附加上下文信息
func WithMessagef(err error, format string, args ...any) error {
	return errors.WithMessagef(err, format, args...)
}

// Unwrap 解包一层错误
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Cause 完整解包错误，得到最内层错误
func Cause(err error) error {
	return errors.Cause(err)
}

// Join 将多个错误包装为一个错误
func Join(errs ...error) error {
	return errors.Join(errs...)
}
