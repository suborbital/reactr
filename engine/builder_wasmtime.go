//go:build wasmtime
// +build wasmtime

package engine

import (
	"github.com/suborbital/reactr/engine/api"
	"github.com/suborbital/reactr/engine/moduleref"
	"github.com/suborbital/reactr/engine/runtime"
	runtimewasmtime "github.com/suborbital/reactr/engine/runtime/wasmtime"
)

func runtimeBuilder(ref *moduleref.WasmModuleRef) runtime.RuntimeBuilder {
	return runtimewasmtime.NewBuilder(ref, api.API()...)
}
