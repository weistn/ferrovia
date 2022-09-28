package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
)

type ErrorCode int

type Error struct {
	code     ErrorCode
	location errlog.LocationRange
	args     []string
}

func NewError(code ErrorCode, loc errlog.LocationRange, args ...string) *Error {
	return &Error{code: code, location: loc, args: args}
}

// Error ...
func (e *Error) Error() string {
	return e.ToString()
}

func (e *Error) ToString() string {
	return "TODO"
}
