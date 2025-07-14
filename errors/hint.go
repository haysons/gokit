package errors

import (
	"github.com/cockroachdb/errors"
)

// WithHint 为错误附加提示信息，提示信息一般用于直接展示给用户
func WithHint(err error, hint string) error {
	return errors.WithHint(err, hint)
}

// WithHintf 为错误附加格式化提示信息，提示信息一般用于直接展示给用户
func WithHintf(err error, format string, args ...any) error {
	return errors.WithHintf(err, format, args...)
}

// GetHint 获取错误中包含的最后一条提示信息
func GetHint(err error) string {
	hints := GetAllHints(err)
	if len(hints) == 0 {
		return ""
	}
	return hints[len(hints)-1]
}

// GetAllHints 获取错误中包含的全部提示信息
func GetAllHints(err error) []string {
	return errors.GetAllHints(err)
}

// FlattenHints 将全部提示信息扁平化为一条提示信息返回
func FlattenHints(err error) string {
	return errors.FlattenHints(err)
}
