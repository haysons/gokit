package errors

import (
	"context"

	"github.com/cockroachdb/errors"
)

// EncodeError 将err进行编码，编码后为一个基于protobuf生成的结构体，主要用于跨网络进行错误传输
func EncodeError(ctx context.Context, err error) errors.EncodedError {
	return errors.EncodeError(ctx, err)
}

// DecodeError 解码错误，还原错误类型
func DecodeError(ctx context.Context, enc errors.EncodedError) error {
	return errors.DecodeError(ctx, enc)
}
