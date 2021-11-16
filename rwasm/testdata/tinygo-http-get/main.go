package main

import (
	"go.suborbital.dev/runnable"
)

type TinygoHttpGet struct{}

func (h TinygoHttpGet) Run(input []byte) ([]byte, error) {
	headers := map[string]string{}
	headers["foo"] = "bar"

	res, err := runnable.POST(string(input), []byte("foobar"), headers)
	if err != nil {
		return nil, err
	}

	runnable.Info(string(res))

	return res, nil
}

// initialize runnable, do not edit //
func main() {
	runnable.Use(TinygoHttpGet{})
}
