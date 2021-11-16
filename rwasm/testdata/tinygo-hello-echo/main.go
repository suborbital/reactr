package main

import (
	"go.suborbital.dev/runnable"
)

type Hello struct{}

func (h Hello) Run(input []byte) ([]byte, error) {
	return []byte("Hello, " + string(input)), nil
}

func main() {
	runnable.Use(Hello{})
}
