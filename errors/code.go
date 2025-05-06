package errors

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
	"github.com/cockroachdb/errors/extgrpc"
	"github.com/cockroachdb/errors/exthttp"
	"github.com/cockroachdb/errors/markers"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/codes"
)

// NewWithCode 创建带有业务状态码的错误
func NewWithCode(code int, msg string) error {
	return &withCode{cause: errors.NewWithDepth(1, msg), code: code}
}

// NewWithCodef 使用格式化方式创建带有业务状态码的错误
func NewWithCodef(code int, format string, args ...any) error {
	return &withCode{cause: errors.NewWithDepthf(1, format, args...), code: code}
}

// WithCode 为错误附加业务状态码，业务状态码相同时，则直接认为是同一个错误
func WithCode(err error, code int) error {
	if err == nil {
		return nil
	}
	return &withCode{cause: err, code: code}
}

// GetCode 获取错误中包含的业务状态码
func GetCode(err error) int {
	if err == nil {
		return 0
	}
	if v, ok := markers.If(err, func(err error) (any, bool) {
		if w, ok := err.(*withCode); ok {
			return w.code, true
		}
		return nil, false
	}); ok {
		return v.(int)
	}
	return 0
}

// WithGrpcCode 为错误附加grpc状态码
func WithGrpcCode(err error, code codes.Code) error {
	return extgrpc.WrapWithGrpcCode(err, code)
}

// GetGrpcCode 获取错误中的grpc状态码
func GetGrpcCode(err error) codes.Code {
	return extgrpc.GetGrpcCode(err)
}

// WithHttpCode 为错误附加http状态码
func WithHttpCode(err error, code int) error {
	return exthttp.WrapWithHTTPCode(err, code)
}

// GetHttpCode 获取错误中的业务状态码
func GetHttpCode(err error, defaultCode int) int {
	return exthttp.GetHTTPCode(err, defaultCode)
}

// withCode 使用业务状态码包装一个错误
type withCode struct {
	cause error
	code  int
}

func (w *withCode) Error() string { return fmt.Sprintf("code=%d, %v", w.code, w.cause) }

func (w *withCode) Cause() error { return w.cause }

func (w *withCode) Unwrap() error { return w.cause }

func (w *withCode) Format(s fmt.State, verb rune) { errors.FormatError(w, s, verb) }

func encodeWithCode(_ context.Context, err error) (string, []string, proto.Message) {
	w := err.(*withCode)
	details := []string{fmt.Sprintf("code: %d", w.code)}
	payload := &exthttp.EncodedHTTPCode{Code: uint32(w.code)}
	return "", details, payload
}

func decodeWithCode(_ context.Context, cause error, _ string, _ []string, payload proto.Message) error {
	wp := payload.(*exthttp.EncodedHTTPCode)
	return &withCode{cause: cause, code: int(wp.Code)}
}

func init() {
	errbase.RegisterWrapperEncoder(errbase.GetTypeKey((*withCode)(nil)), encodeWithCode)
	errbase.RegisterWrapperDecoder(errbase.GetTypeKey((*withCode)(nil)), decodeWithCode)
}
