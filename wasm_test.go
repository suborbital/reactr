package hive

import (
	"fmt"
	"testing"
)

func TestWasmRunner(t *testing.T) {
	h := New()

	doWasm := h.Handle("wasm", NewWasm("./wasm/wasm_runner_bg.wasm"))

	grp := NewGroup()
	for i := 0; i < 50000; i++ {
		grp.Add(doWasm(fmt.Sprintf("world %d", i)))
	}

	if err := grp.Wait(); err != nil {
		t.Error(err)
	}
}
