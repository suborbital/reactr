//go:build tinygo.wasm

package errors

import (
	"errors"
)

// RunErr adds a status code for use in FFI calls to return_result()
type RunErr struct {
	error
	Code int
}

type HostErr error

func NewError(code int, message string) RunErr {
	return RunErr{errors.New(message), code}
}

func NewHostError(message string) HostErr {
	return errors.New(message).(HostErr)
}

// WithCode creates a new RunErr from an existing error and a status code
func WithCode(err error, code int) RunErr {
	return RunErr{err, code}
}
