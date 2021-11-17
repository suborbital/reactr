package main

import "github.com/suborbital/reactr/api/tinygo/runnable"

type TinygoLog struct{}

func (h TinygoLog) Run(input []byte) ([]byte, error) {
	runnable.Info(string(input))
	runnable.Info("info log")
	runnable.Error("some error")

	warnMsg := "warning message"
	runnable.Warnf("some %s", warnMsg)

	runnable.Debug("debug message")

	return []byte(""), nil
}

// initialize runnable, do not edit //
func main() {
	runnable.Use(TinygoLog{})
}
