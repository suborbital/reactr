package main

import "github.com/suborbital/reactr/api/tinygo/runnable"

type Cache struct{}

func (h Cache) Run(input []byte) ([]byte, error) {
	runnable.CacheSet("name", string(input), 0)

	return runnable.CacheGet("name")
}

func main() {
	runnable.Use(Cache{})
}
