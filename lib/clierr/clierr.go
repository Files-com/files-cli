// package clierr implements an error type that includes an error code for
// in order to ensure the CLI can return the correct exit code to the user.
package clierr

import (
	"errors"
	"fmt"
)

// ErrorCode is an integer type representing different error codes.
type ErrorCode int

// emulating rclone error codes from https://rclone.org/docs/#list-of-exit-codes
// but leaving out the ones that are not relevant to this CLI (8, 9, 10)
const (
	// ErrorCodeDefault is the default error code used when none is provided.
	ErrorCodeDefault ErrorCode = iota + 1
	// ErrorCodeUsage is the error code for usage errors.
	ErrorCodeUsage
	// ErrorCodeDirectoryNotFound is the error code for directory not found errors.
	ErrorCodeDirectoryNotFound
	// ErrorCodeFileNotFound is the error code for file not found errors.
	ErrorCodeFileNotFound
	// ErrorCodeTemporary is the error code for temporary errors.
	ErrorCodeTemporary
	// ErrorCodeNotRetryable is the error code for not retryable errors.
	ErrorCodeNotRetryable
	// ErrorCodeFatal is the error code for fatal errors.
	ErrorCodeFatal
)

// CliError represents an error with an associated status code.
type CliError struct {
	Code ErrorCode
	Err  error
}

// Error implements the error interface for StatusError.
func (e *CliError) Error() string {
	return fmt.Sprintf("%s - status (%d)", e.Err.Error(), e.Code)
}

// Unwrap implements the Unwrap method for the error interface.
func (e *CliError) Unwrap() error {
	return e.Err
}

// New creates a new CliError with the given code and error.
func New(code ErrorCode, err error) *CliError {
	return &CliError{
		Code: code,
		Err:  err,
	}
}

// Error creates a new CliError with the given code and message.
func Error(code ErrorCode, msg string) *CliError {
	return &CliError{
		Code: code,
		Err:  errors.New(msg),
	}
}

// Errorf creates a new CliError with the given code and formatted message.
func Errorf(code ErrorCode, format string, args ...any) *CliError {
	return &CliError{
		Code: code,
		Err:  fmt.Errorf(format, args...),
	}
}

// From checks if the error is a CliError and returns it.
// If not, it creates a new CliError with the default code.
func From(err error) *CliError {
	if cliErr, ok := err.(*CliError); ok {
		return cliErr
	}
	return New(ErrorCodeDefault, err)
}
