package runnable

// The Runnable interface is all that needs to be implemented by a Reactr runnable.
type Runnable interface {
	Run(input []byte) ([]byte, error)
}

func Use(runnable Runnable) {
	runnable_ = runnable
}
