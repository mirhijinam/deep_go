package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	headMsg := fmt.Sprintf("%d errors occured:\n", len(e.errs))

	errMsgs := make([]string, 0)
	for _, err := range e.errs {
		errMsgs = append(errMsgs, fmt.Sprintf("* %s", err.Error()))
	}

	return headMsg + "\t" + strings.Join(errMsgs, "\t") + "\n"
}

func Append(err error, errs ...error) *MultiError {
	multiErr, ok := err.(*MultiError)
	if !ok {
		return &MultiError{errs}
	}

	multiErr.errs = append(multiErr.errs, errs...)

	return multiErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
