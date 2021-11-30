package wasmtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
)

func TestPythonHelloWorld(t *testing.T) {
	r := rt.New()

	doPython := r.Register("python", rwasm.NewRunner("../../api/python/python.wasm"))

	start := time.Now()

	_, err := doPython(`log_info("hello from python!")`).Then()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(time.Since(start).Milliseconds())
}
