package main

import (
	"fmt"
	"log"

	"github.com/suborbital/hive"
)

func main() {
	h := hive.New()

	doWasm := h.Handle("wasm", hive.NewWasm("./pkg/wasm_test_bg.wasm"))

	grp := hive.NewGroup()
	for i := 0; i < 50000; i++ {
		grp.Add(doWasm(fmt.Sprintf("world %d", i)))
	}

	if err := grp.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("done")
}
