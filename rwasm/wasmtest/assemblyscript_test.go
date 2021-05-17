package wasmtest

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
)

func TestASEcho(t *testing.T) {
	r := rt.New()

	// test a WASM module that is loaded directly instead of through the bundle
	doWasm := r.Register("as-echo", rwasm.NewRunner("../testdata/as-echo/as-echo.wasm"))

	res, err := doWasm("from AssemblyScript!").Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	fmt.Println(string(res.([]byte)))

	if string(res.([]byte)) != "hello, from AssemblyScript!" {
		t.Error("as-echo failed, got:", string(res.([]byte)))
	}
}

func TestASFetch(t *testing.T) {
	r := rt.New()

	// test a WASM module that is loaded directly instead of through the bundle
	doWasm := r.Register("as-fetch", rwasm.NewRunner("../testdata/as-fetch/as-fetch.wasm"))

	res, err := doWasm("https://1password.com").Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)[:100]) != "<!doctype html><html lang=en data-language-url=/><head><meta charset=utf-8><meta name=viewport conte" {
		t.Error("as-fetch failed, got:", string(res.([]byte)[:100]))
	}
}

func TestASLargeData(t *testing.T) {
	r := rt.New()

	// test a WASM module that is loaded directly instead of through the bundle
	doWasm := r.Register("as-echo", rwasm.NewRunner("../testdata/as-echo/as-echo.wasm"))

	res, err := doWasm(largeInput).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)) != "hello, "+largeInput {
		t.Error("as-test failed, got:", string(res.([]byte)))
	}
}
