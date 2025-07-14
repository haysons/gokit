package errors

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"net/http"
	"testing"

	cerrors "github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewAndNewf(t *testing.T) {
	err1 := New("simple error")
	assert.NotNil(t, err1)
	assert.Equal(t, "simple error", err1.Error())

	err2 := Newf("formatted %s: %d", "error", 42)
	assert.NotNil(t, err2)
	assert.Equal(t, "formatted error: 42", err2.Error())
}

func TestWrapAndWrapf_Unwrap_Cause(t *testing.T) {
	orig := errors.New("root cause")

	w1 := Wrap(orig, "context")
	assert.NotNil(t, w1)
	assert.Equal(t, "context: root cause", w1.Error())
	assert.Equal(t, orig, Cause(Unwrap(w1)))

	w2 := Wrapf(w1, "more %s", "info")
	assert.NotNil(t, w2)
	assert.Equal(t, "more info: context: root cause", w2.Error())
	assert.Equal(t, orig, Cause(w2))
}

func TestWithMessage(t *testing.T) {
	origErr := errors.New("original error")
	wrapped := WithMessage(origErr, "context added")

	assert.Error(t, wrapped)
	assert.EqualError(t, wrapped, "context added: original error")
	assert.True(t, errors.Is(wrapped, origErr), "wrapped error should match original with errors.Is")
}

func TestWithMessagef(t *testing.T) {
	origErr := errors.New("read failed")
	wrapped := WithMessagef(origErr, "operation %s on file %s", "read", "/tmp/data.txt")

	assert.Error(t, wrapped)
	assert.EqualError(t, wrapped, "operation read on file /tmp/data.txt: read failed")
	assert.True(t, errors.Is(wrapped, origErr), "wrapped error should match original with errors.Is")
}

func TestIs(t *testing.T) {
	errA := New("foo")
	errB := New("bar")
	errC := New("bar")
	wrappedA := Wrap(errA, "ctx")

	assert.True(t, Is(errA, errA), "same error should be Is")
	assert.False(t, Is(errA, errB), "different errors not Is")
	assert.True(t, Is(wrappedA, errA), "wrapped contains target")
	assert.False(t, Is(wrappedA, errB), "wrapped does not contain non-target")
	assert.True(t, Is(errB, errC))
}

type MyErr struct{ msg string }

func (e *MyErr) Error() string { return e.msg }

func TestAs(t *testing.T) {
	orig := &MyErr{"my error"}
	wrapped := Wrap(orig, "ctx")

	var got *MyErr
	ok := As(wrapped, &got)
	assert.True(t, ok, "wrapped error should As to *MyErr")
	assert.Equal(t, orig, got)
}

func TestJoin(t *testing.T) {
	e1 := New("one")
	e2 := New("two")
	joined := Join(e1, e2)

	// errors.Is should find both
	assert.True(t, Is(joined, e1))
	assert.True(t, Is(joined, e2))
}

func TestFormat(t *testing.T) {
	assert.Equal(t, "", Format(nil))
	err := demo2()
	errFormatted := Format(err)
	t.Log(errFormatted)
}

func TestIntegrationWithCockroachErrors(t *testing.T) {
	// ensure our Wrap/New behave like cockroachdb/errors
	orig := cerrors.New("orig")
	wrapped := Wrap(orig, "ctx")
	assert.True(t, cerrors.Is(wrapped, orig))
	// roundtrip encode/decode
	enc := cerrors.EncodeError(context.Background(), wrapped)
	dec := cerrors.DecodeError(context.Background(), enc)
	assert.True(t, cerrors.Is(dec, orig))
}

func TestNewBiz(t *testing.T) {
	code := 1001
	hint := "请检查输入参数"
	msg := "参数无效"

	err := NewBiz(code, hint, msg)
	assert.Error(t, err)

	// 检查 error 消息是否匹配
	assert.Contains(t, err.Error(), msg)

	// 检查附加信息
	assert.Equal(t, code, GetCode(err))
	assert.Equal(t, http.StatusOK, GetHttpCode(err, 0))
	assert.Equal(t, codes.FailedPrecondition, GetGrpcCode(err))
	assert.Equal(t, hint, GetAllHints(err)[0])
}

func TestNewBizf(t *testing.T) {
	code := 1002
	hint := "用户名不能为空"
	format := "字段 %s 无效: %s"
	field := "username"
	reason := "不能为空"

	err := NewBizf(code, hint, format, field, reason)
	assert.Error(t, err)

	expectedMsg := "字段 username 无效: 不能为空"
	assert.Contains(t, err.Error(), expectedMsg)

	assert.Equal(t, code, GetCode(err))
	assert.Equal(t, http.StatusOK, GetHttpCode(err, 0))
	assert.Equal(t, codes.FailedPrecondition, GetGrpcCode(err))
	assert.Equal(t, hint, GetAllHints(err)[0])
}
