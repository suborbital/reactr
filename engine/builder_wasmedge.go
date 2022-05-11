//go:build wasmedge
// +build wasmedge

package engine

import (
	"github.com/suborbital/reactr/engine/api"
	"github.com/suborbital/reactr/engine/moduleref"
	"github.com/suborbital/reactr/engine/runtime"
	runtimewasmedge "github.com/suborbital/reactr/engine/runtime/wasmedge"
)

func runtimeBuilder(ref *moduleref.WasmModuleRef, api api.HostAPI) runtime.RuntimeBuilder {
	return runtimewasmedge.NewBuilder(ref, api)
}
