//go:build !wasmtime
// +build !wasmtime

// these tests are skipped with the wasmtime runtime, since it did not inherit the legacy baggage
package wasmtest

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/engine"
	"github.com/suborbital/reactr/scheduler"
)

func TestWasmLegacyInit(t *testing.T) {
	e := engine.New()

	e.RegisterFromFile("legacy", "../testdata/legacy/legacy.wasm")

	job := scheduler.NewJob("legacy", "Connor")

	res, err := e.Do(job).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)) != "hello Connor" {
		t.Error(fmt.Errorf("hello Connor, got %s", string(res.([]byte))))
	}
}
