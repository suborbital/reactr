package runnable

import "errors"

// RunErr adds a status code for use in FFI calls to return_result()
type RunErr struct {
	error
	Code int
}

// NewError creates a new RunErr, which is a normal Go error plus an error code
func NewError(code int, message string) RunErr {
	return RunErr{errors.New(message), code}
}

type HostErr error

func NewHostError(message string) HostErr {
	return errors.New(message).(HostErr)
}
