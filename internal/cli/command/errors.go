package command

import "fmt"

// Error represents a command execution error with optional usage flag.
type Error struct {
	Message string
	Usage   bool
}

func (e *Error) Error() string {
	return e.Message
}

// IsUsage returns true if the error was caused by incorrect usage.
func (e *Error) IsUsage() bool {
	return e.Usage
}

// NewError creates a new Error with the given message and usage flag.
func NewError(msg string, usage bool) *Error {
	return &Error{Message: msg, Usage: usage}
}

// NewUsageError creates a new Error indicating incorrect command usage.
func NewUsageError(msg string) *Error {
	return &Error{Message: msg, Usage: true}
}

// NewRuntimeError creates a new Error for runtime failures.
func NewRuntimeError(msg string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(msg, args...), Usage: false}
}
