package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithStackAndGetStackTrace(t *testing.T) {
	baseErr := New("something failed")
	stackErr := WithStack(baseErr)

	assert.NotEqual(t, baseErr, stackErr)

	trace := GetStackTrace(stackErr)
	assert.NotNil(t, trace)
	assert.Greater(t, len(trace), 0, "stack trace should not be empty")

	wrapErr := demo2()
	trace = GetStackTrace(wrapErr)
	t.Logf("source: %v", GetSource(wrapErr))
	t.Logf("%+v\n\n", trace)

	codeErr := NewWithCode(1001, "something failed")
	wrapErr = Wrap(codeErr, "wrap code")
	t.Logf("source: %v", GetSource(wrapErr))
	trace = GetStackTrace(wrapErr)
	t.Logf("%+v", trace)
}

func demo2() error {
	return demo1()
}

func demo1() error {
	baseErr := New("demo1")
	wrapErr := Wrap(baseErr, "wrap")
	return wrapErr
}
