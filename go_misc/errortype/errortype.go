package errortype

import (
	"fmt"
	"strings"
)

var StackedPrefix = "\n<-- "

type StackedError struct {
	Errors []error `json:"error"` // earliest = root cause
}

func NewStackedError(err error) *StackedError {
	return &StackedError{[]error{err}}
}

func StackedErrorf(formatstr string, params ...interface{}) *StackedError {
	return NewStackedError(fmt.Errorf(formatstr, params...))
}

func (err *StackedError) Stack(newerr error) *StackedError {
	err.Errors = append(err.Errors, newerr)
	return err
}

func (err *StackedError) StackStrf(errstring string, params ...interface{}) *StackedError {
	err.Errors = append(err.Errors, fmt.Errorf(errstring, params...))
	return err
}

func (err *StackedError) Error() string {
	errStrList := make([]string, len(err.Errors))
	for i, childerr := range err.Errors {
		errStrList[i] = childerr.Error()
	}
	return strings.Join(errStrList, StackedPrefix)
}
