package main

import (
	"go.suborbital.dev/runnable"
)

type Cache struct{}

func (h Cache) Run(input []byte) ([]byte, error) {
	runnable.CacheSet(string(input), "hello world", 0)

	return runnable.CacheGet(string(input))
}

func main() {
	runnable.Use(Cache{})
}
