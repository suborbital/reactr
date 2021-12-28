package runnable

// The Runnable interface is all that needs to be implemented by a Reactr runnable.
type Runnable interface {
	Run(input []byte) ([]byte, error)
}

// RunErr adds a status code for use in FFI calls to return_result()
type RunErr struct {
	error
	Code int
}

type HostErr error
