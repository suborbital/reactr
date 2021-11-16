package main

import (
	"go.suborbital.dev/runnable"
)

type TinygoReq struct{}

func (h TinygoReq) Run(input []byte) ([]byte, error) {
	method := runnable.Method()
	url := runnable.URL()

	param := runnable.URLParam("foobar")

	runnable.Infof("%s: %s?%s", method, url, param)
	return []byte("Success"), nil
}

// initialize runnable, do not edit //
func main() {
	runnable.Use(TinygoReq{})
}
