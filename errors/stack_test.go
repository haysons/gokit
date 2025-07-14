package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithStackAndGetStackTrace(t *testing.T) {
	baseErr := New("something failed")
	stackErr := WithStack(baseErr)

	assert.NotEqual(t, baseErr, stackErr)
	assert.True(t, errors.Is(stackErr, baseErr))

	trace := GetStackTrace(stackErr)
	assert.NotNil(t, trace)
	assert.Greater(t, len(trace), 0, "stack trace should not be empty")

	wrapErr := demo2()
	stack := GetStack(wrapErr)
	stack = stack[:3]
	assert.Equal(t, []string{"stack_test.go:35 demo1", "stack_test.go:30 demo2", "stack_test.go:21 TestWithStackAndGetStackTrace"}, stack)
	stack = GetStack(errors.New("something failed"))
	assert.Len(t, stack, 0)
}

func demo2() error {
	return demo1()
}

func demo1() error {
	baseErr := New("demo1")
	wrapErr := Wrap(baseErr, "wrap")
	return wrapErr
}
