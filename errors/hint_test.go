package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHints(t *testing.T) {
	baseErr := New("something went wrong")
	hintedErr := WithHint(baseErr, "try restarting the application")

	hints := GetAllHints(hintedErr)
	assert.Len(t, hints, 1)
	assert.Equal(t, "try restarting the application", hints[0])

	flat := FlattenHints(hintedErr)
	assert.Equal(t, "try restarting the application", flat)
}

func TestHintf(t *testing.T) {
	baseErr := New("file not found")
	hintedErr := WithHintf(baseErr, "check if file %s exists", "/tmp/test.txt")

	hints := GetAllHints(hintedErr)
	assert.Len(t, hints, 1)
	assert.Equal(t, "check if file /tmp/test.txt exists", hints[0])

	flat := FlattenHints(hintedErr)
	assert.Equal(t, "check if file /tmp/test.txt exists", flat)
}

func TestMultipleHints(t *testing.T) {
	err := New("database error")
	err = WithHint(err, "check your credentials")
	err = WithHint(err, "ensure the DB is reachable")

	hints := GetAllHints(err)
	assert.Len(t, hints, 2)
	assert.Equal(t, "check your credentials", hints[0])
	assert.Equal(t, "ensure the DB is reachable", hints[1])

	flat := FlattenHints(err)
	assert.Contains(t, flat, "check your credentials")
	assert.Contains(t, flat, "ensure the DB is reachable")

	hint := GetHint(err)
	assert.Equal(t, "ensure the DB is reachable", hint)

	err = WithHint(err, "check your address")
	hint = GetHint(err)
	assert.Equal(t, "check your address", hint)

	hint = GetHint(errors.New("some other error"))
	assert.Equal(t, "", hint)
}
