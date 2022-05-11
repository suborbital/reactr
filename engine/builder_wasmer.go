//go:build wasmer
// +build wasmer

package engine

import (
	"github.com/suborbital/reactr/engine/api"
	"github.com/suborbital/reactr/engine/moduleref"
	"github.com/suborbital/reactr/engine/runtime"
	runtimewasmer "github.com/suborbital/reactr/engine/runtime/wasmer"
)

func runtimeBuilder(ref *moduleref.WasmModuleRef, api api.HostAPI) runtime.RuntimeBuilder {
	return runtimewasmer.NewBuilder(ref, api)
}
