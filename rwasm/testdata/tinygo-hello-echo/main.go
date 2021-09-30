package main

package main

import (
	"fmt"

	"github.com/jagger27/hello-wasm/suborbital"
)

type Hello struct{}

func (h Hello) Run(input []byte) []byte {
	return []byte("Hello, " + string(input))
}

// insert here
func main() {
	suborbital.RUNNABLE = Hello{}
}
