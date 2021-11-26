//go:build wasmedge
// +build wasmedge

package rwasm

import (
	"github.com/suborbital/reactr/rwasm/api"
	"github.com/suborbital/reactr/rwasm/moduleref"
	"github.com/suborbital/reactr/rwasm/runtime"
	runtimewasmedge "github.com/suborbital/reactr/rwasm/runtime/wasmedge"
)

func runtimeBuilder(ref *moduleref.WasmModuleRef) runtime.RuntimeBuilder {
	return runtimewasmedge.NewBuilder(ref, api.API()...)
}
