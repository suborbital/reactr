package main

import (
	"github.com/suborbital/libtinygo"
)

type Hello struct{}

func (h Hello) Run(input []byte) []byte {
	return []byte("Hello, " + string(input))
}

// insert here
func main() {
	libtinygo.RUNNABLE = Hello{}
}
