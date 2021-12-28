package runnable

import "errors"

// NewError creates a new RunErr, which is a normal Go error plus an error code
func NewError(code int, message string) RunErr {
	return RunErr{errors.New(message), code}
}

func NewHostError(message string) HostErr {
	return errors.New(message).(HostErr)
}
