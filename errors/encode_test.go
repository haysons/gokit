package errors

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeError(t *testing.T) {
	ctx := context.Background()

	original := WithCode(New("encode decode test"), 1001)

	encoded := EncodeError(ctx, original)
	assert.NotNil(t, encoded, "EncodedError should not be nil")

	decoded := DecodeError(ctx, encoded)
	assert.NotNil(t, decoded, "Decoded error should not be nil")

	assert.Contains(t, decoded.Error(), "encode decode test")
	assert.Equal(t, GetCode(decoded), 1001)

	original = New("org")
	encoded = EncodeError(ctx, original)
	decoded = DecodeError(ctx, encoded)
	assert.True(t, errors.Is(original, decoded), "Decoded error should be identified as same by original")
}
