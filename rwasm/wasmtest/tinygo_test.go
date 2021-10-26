package wasmtest

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
)

func TestWasmRunnerTinyGo(t *testing.T) {
	r := rt.New()

	// test a WASM module that is loaded directly instead of through the bundle
	doWasm := r.Register("wasm", rwasm.NewRunner("../testdata/tinygo-hello-echo/tinygo-hello-echo.wasm"))

	res, err := doWasm("world").Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)) != "Hello, world" {
		t.Errorf("expected Hello world got %q", string(res.([]byte)))
	}
}
