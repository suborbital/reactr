package wasmtest

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rt"
)

func TestWasmCacheGetSetRustToSwift(t *testing.T) {
	setJob := rt.NewJob("rust-set", "very important")
	getJob := rt.NewJob("swift-get", "")

	_, err := sharedRT.Do(setJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set cache value"))
		return
	}

	r2, err := sharedRT.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to get cache value"))
		return
	}

	if string(r2.([]byte)) != "very important" {
		t.Error(fmt.Errorf("did not get expected output"))
	}
}

func TestWasmCacheGetSetSwiftToRust(t *testing.T) {
	setJob := rt.NewJob("swift-set", "very important")
	getJob := rt.NewJob("rust-get", "")

	_, err := sharedRT.Do(setJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set cache value"))
		return
	}

	r2, err := sharedRT.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to get cache value"))
		return
	}

	if string(r2.([]byte)) != "very important" {
		t.Error(fmt.Errorf("did not get expected output"))
	}
}

func TestWasmCacheGetSetSwiftToAS(t *testing.T) {
	setJob := rt.NewJob("swift-set", "very important")
	getJob := rt.NewJob("as-get", "")

	_, err := sharedRT.Do(setJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set cache value"))
		return
	}

	r2, err := sharedRT.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to get cache value"))
		return
	}

	if string(r2.([]byte)) != "very important" {
		t.Error(fmt.Errorf("did not get expected output"))
	}
}

func TestWasmCacheGetSetASToRust(t *testing.T) {
	setJob := rt.NewJob("as-set", "very important")
	getJob := rt.NewJob("rust-get", "")

	_, err := sharedRT.Do(setJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set cache value"))
		return
	}

	r2, err := sharedRT.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to get cache value"))
		return
	}

	if string(r2.([]byte)) != "very important" {
		t.Error(fmt.Errorf("did not get expected output"))
	}
}
