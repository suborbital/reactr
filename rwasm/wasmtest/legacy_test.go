//go:build !wasmtime

// these tests are skipped with the wasmtime runtime, since it did not inherit the legacy baggage
package wasmtest

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
)

func TestWasmLegacyInit(t *testing.T) {
	r := rt.New()
	r.Register("legacy", rwasm.NewRunner("../testdata/legacy/legacy.wasm"))

	job := rt.NewJob("legacy", "Connor")

	res, err := r.Do(job).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)) != "hello Connor" {
		t.Error(fmt.Errorf("hello Connor, got %s", string(res.([]byte))))
	}
}
