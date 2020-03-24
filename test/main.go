package main

import (
	"fmt"
	"log"

	"github.com/suborbital/hive"
)

func main() {
	h := hive.New()

	doWasm := h.Handle("wasm", hive.NewWasm("./pkg/wasm_test_bg.wasm"))

	res, err := doWasm("world").Then()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.(string))

	fmt.Println("done")
}
