package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (m *MultiError) Unwrap() error {
	if m == nil || len(m.errors) == 0 {
		return nil
	}

	errs := m.errors[1:]
	if len(errs) == 1 {
		return errs[0]
	}

	return &MultiError{
		errors: errs,
	}
}

func (m *MultiError) Error() string {
	length := len(m.errors)
	if length == 0 {
		return "no errors"
	}

	if length == 1 {
		return m.errors[0].Error()
	}

	var builder strings.Builder
	builder.WriteString(strconv.FormatInt(int64(length), 10))
	builder.WriteString(" errors occurred:\n")
	for _, err := range m.errors {
		builder.WriteString("\t* ")
		builder.WriteString(err.Error())
	}
	builder.WriteString("\n")

	return builder.String()
}

func (m *MultiError) Is(target error) bool {
	if m == nil {
		return m == target
	}

	for _, err := range m.errors {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func (m *MultiError) As(target any) bool {
	if m == nil {
		return m == target
	}

	for _, err := range m.errors {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}

func Append(err error, errs ...error) *MultiError {
	if err == nil && len(errs) == 0 {
		return nil
	}

	if err == nil {
		return &MultiError{
			errors: errs,
		}
	}

	mErr, ok := err.(*MultiError)
	if !ok {
		mErr = &MultiError{
			errors: make([]error, 0, len(errs)+1),
		}
		mErr.errors = append(mErr.errors, err)
	}

	mErr.errors = append(mErr.errors, errs...)

	return mErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occurred:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiError_Unwrap(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")

	merr := Append(nil, err1, err2, err3)

	var current error = merr
	assert.NotNil(t, current, "current error should not be nil")
	assert.Equal(t, "3 errors occurred:\n\t* err1\t* err2\t* err3\n", current.Error())

	current = errors.Unwrap(current)
	assert.NotNil(t, current, "current error should not be nil")
	assert.Equal(t, "2 errors occurred:\n\t* err2\t* err3\n", current.Error())

	current = errors.Unwrap(current)
	assert.NotNil(t, current, "current error should not be nil")
	assert.Equal(t, "err3", current.Error())

	current = errors.Unwrap(current)
	assert.Nil(t, current, "should be nil at the end")
}

func TestMultiError_Is(t *testing.T) {
	target := errors.New("target error")
	other := errors.New("other error")

	err := Append(nil, other, target)

	assert.True(t, errors.Is(err, target), "errors.Is should find target error")
	assert.False(t, errors.Is(err, errors.New("nonexistent")), "errors.Is should not find non-existent error")
}

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func TestMultiError_As(t *testing.T) {

	custom := &customError{"custom error"}
	err := Append(nil, errors.New("a"), custom, errors.New("b"))

	var out *customError
	assert.True(t, errors.As(err, &out), "errors.As should find customError")
	assert.Equal(t, "custom error", out.msg)
}
