package main

import (
	"github.com/suborbital/runnable"
)

type Hello struct{}

func (h Hello) Run(input []byte) ([]byte, error) {
	return []byte("Hello, " + string(input)), nil
}

// insert here
func main() {
	runnable.Use(Hello{})
}
