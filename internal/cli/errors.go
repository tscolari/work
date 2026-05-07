package cli

import "fmt"

// Error carries a process exit code along with a user-facing message.
//
// Code 1 = user error (bad input, missing prereqs).
// Code 2 = unexpected failure.
type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string { return e.Msg }

func userErr(format string, args ...any) *Error {
	return &Error{Code: 1, Msg: fmt.Sprintf(format, args...)}
}

func sysErr(format string, args ...any) *Error {
	return &Error{Code: 2, Msg: fmt.Sprintf(format, args...)}
}
