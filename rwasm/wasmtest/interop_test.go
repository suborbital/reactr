package wasmtest

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
)

func TestWasmCacheGetSetRustToSwift(t *testing.T) {
	r := rt.New()
	r.Register("rust-set", rwasm.NewRunner("../testdata/rust-set/rust-set.wasm"))
	r.Register("swift-get", rwasm.NewRunner("../testdata/swift-get/swift-get.wasm"))

	setJob := rt.NewJob("rust-set", "very important")
	getJob := rt.NewJob("swift-get", "")

	_, err := r.Do(setJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set cache value"))
		return
	}

	r2, err := r.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to get cache value"))
		return
	}

	if string(r2.([]byte)) != "very important" {
		t.Error(fmt.Errorf("did not get expected output"))
	}
}

func TestWasmCacheGetSetSwiftToRust(t *testing.T) {
	r := rt.New()
	r.Register("swift-set", rwasm.NewRunner("../testdata/swift-set/swift-set.wasm"))
	r.Register("rust-get", rwasm.NewRunner("../testdata/rust-get/rust-get.wasm"))

	setJob := rt.NewJob("swift-set", "very important")
	getJob := rt.NewJob("rust-get", "")

	_, err := r.Do(setJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set cache value"))
		return
	}

	r2, err := r.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to get cache value"))
		return
	}

	if string(r2.([]byte)) != "very important" {
		t.Error(fmt.Errorf("did not get expected output"))
	}
}

func TestWasmCacheGetSetSwiftToAS(t *testing.T) {
	r := rt.New()
	r.Register("swift-set", rwasm.NewRunner("../testdata/swift-set/swift-set.wasm"))
	r.Register("as-get", rwasm.NewRunner("../testdata/as-get/as-get.wasm"))

	setJob := rt.NewJob("swift-set", "very important")
	getJob := rt.NewJob("as-get", "")

	_, err := r.Do(setJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set cache value"))
		return
	}

	r2, err := r.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to get cache value"))
		return
	}

	if string(r2.([]byte)) != "very important" {
		t.Error(fmt.Errorf("did not get expected output"))
	}
}

func TestWasmCacheGetSetASToRust(t *testing.T) {
	r := rt.New()
	r.Register("as-set", rwasm.NewRunner("../testdata/as-set/as-set.wasm"))
	r.Register("rust-get", rwasm.NewRunner("../testdata/rust-get/rust-get.wasm"))

	setJob := rt.NewJob("as-set", "very important")
	getJob := rt.NewJob("rust-get", "")

	_, err := r.Do(setJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set cache value"))
		return
	}

	r2, err := r.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to get cache value"))
		return
	}

	if string(r2.([]byte)) != "very important" {
		t.Error(fmt.Errorf("did not get expected output"))
	}
}

func TestWasmCacheGetSetTinyGoToRust(t *testing.T) {
	r := rt.New()
	r.Register("tinygo-cache", rwasm.NewRunner("../testdata/tinygo-cache/tinygo-cache.wasm"))
	r.Register("rust-get", rwasm.NewRunner("../testdata/rust-get/rust-get.wasm"))

	setJob := rt.NewJob("tinygo-cache", "very important")
	getJob := rt.NewJob("rust-get", "")

	_, err := r.Do(setJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set cache value"))
		return
	}

	r2, err := r.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to get cache value"))
		return
	}

	if string(r2.([]byte)) != "very important" {
		t.Error(fmt.Errorf("did not get expected output"))
	}
}
