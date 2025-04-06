package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type multiErrorNode struct {
	err  error
	next error
}

func (n *multiErrorNode) Error() string {
	return n.err.Error()
}

func (n *multiErrorNode) Unwrap() error {
	return n.next
}

func (n *multiErrorNode) Is(target error) bool {
	return errors.Is(n.err, target)
}

func (n *multiErrorNode) As(target interface{}) bool {
	return errors.As(n.err, target)
}

type MultiError struct {
	head *multiErrorNode
	len  int
}

func (e *MultiError) Error() string {
	if e == nil || e.head == nil {
		return ""
	}

	var msgs []string
	for node := e.head; node != nil; node = node.next.(*multiErrorNode) {
		msgs = append(msgs, fmt.Sprintf("* %s", node.err.Error()))
		if node.next == nil {
			break
		}
	}

	if len(msgs) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t%s\n", msgs[0])
	}

	return fmt.Sprintf("%d errors occurred:\n\t%s\n", len(msgs), strings.Join(msgs, "\t"))
}

func (e *MultiError) Unwrap() error {
	if e == nil || e.head == nil {
		return nil
	}
	return e.head
}

func Append(err error, errs ...error) *MultiError {
	var root *multiErrorNode
	var last *multiErrorNode
	length := 0

	appendOne := func(e error) {
		if e == nil {
			return
		}
		node := &multiErrorNode{err: e}
		if root == nil {
			root = node
			last = node
		} else {
			last.next = node
			last = node
		}
		length++
	}

	if err != nil {
		if me, ok := err.(*MultiError); ok && me != nil {
			for node := me.head; node != nil; node = node.next.(*multiErrorNode) {
				appendOne(node.err)
				if node.next == nil {
					break
				}
			}
		} else {
			appendOne(err)
		}
	}

	for _, e := range errs {
		if e == nil {
			continue
		}
		if me, ok := e.(*MultiError); ok && me != nil {
			for node := me.head; node != nil; node = node.next.(*multiErrorNode) {
				appendOne(node.err)
				if node.next == nil {
					break
				}
			}
		} else {
			appendOne(e)
		}
	}

	if length == 0 {
		return nil
	}
	return &MultiError{head: root, len: length}
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
	for _, expected := range []error{err1, err2, err3} {
		current = errors.Unwrap(current)
		assert.Equal(t, expected.Error(), current.Error())
	}

	assert.Nil(t, errors.Unwrap(current), "should be nil at the end")
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
