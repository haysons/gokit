package errors

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestNewWithCode(t *testing.T) {
	err := NewWithCode(1001, "something wrong")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "something wrong")
	assert.Equal(t, 1001, GetCode(err))
}

func TestNewfWithCode(t *testing.T) {
	err := NewWithCodef(2002, "error with value: %d", 42)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "value: 42")
	assert.Equal(t, 2002, GetCode(err))
}

func TestWithCode_Wrap(t *testing.T) {
	base := New("base")
	err := WithCode(base, 3003)
	require.Error(t, err)
	assert.Equal(t, 3003, GetCode(err))
	assert.True(t, errors.Is(err, base))
}

func TestWithCode_Nil(t *testing.T) {
	var err error
	assert.Nil(t, WithCode(err, 999))
}

func TestWithGrpcCode(t *testing.T) {
	base := New("grpc test")
	err := WithGrpcCode(base, codes.PermissionDenied)
	require.Error(t, err)
	assert.Equal(t, codes.PermissionDenied, GetGrpcCode(err))
}

func TestWithHttpCode(t *testing.T) {
	base := New("http test")
	err := WithHttpCode(base, 403)
	require.Error(t, err)
	assert.Equal(t, 403, GetHttpCode(err, 500))
	assert.Equal(t, 500, GetHttpCode(nil, 500))
}

func TestEncodeDecodeWithCode(t *testing.T) {
	orig := NewWithCode(1234, "encoded error")
	encoded := EncodeError(context.Background(), orig)
	require.NotNil(t, encoded)

	decoded := DecodeError(context.Background(), encoded)
	require.Error(t, decoded)

	assert.Equal(t, GetCode(orig), GetCode(decoded))
	assert.Equal(t, "encoded error", Cause(decoded).Error())
	assert.True(t, errors.Is(decoded, orig))
}
