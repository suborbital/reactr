package runtimewasmedge

import (
	"github.com/suborbital/reactr/rwasm/runtime"
	"github.com/second-state/WasmEdge-go/wasmedge"
)

// toWasmEdgeHostFn creates a new host funcion from a generic host fn
func toWasmEdgeHostFn(hostFn runtime.HostFn) func(data interface{}, mem *wasmedge.Memory, params []interface{}) ([]interface{}, wasmedge.Result) {
	return func(data interface{}, mem *wasmedge.Memory, params []interface{}) ([]interface{}, wasmedge.Result) {
		hostResult, hostErr := hostFn.HostFn(params...)
		if hostErr != nil {
			return nil, wasmedge.Result_Fail
		}

		return []interface{}{hostResult}, wasmedge.Result_Success
	}
}


// addHostFns adds a list of host functions to an import object
func addHostFns(imports *wasmedge.ImportObject, fns ...runtime.HostFn) {
	for _, fn := range fns {
		wasmHostFn := toWasmEdgeHostFn(fn)

		argsType := make([]wasmedge.ValType, fn.ArgCount)
		for i := 0; i < fn.ArgCount; i++ {
			argsType[i] = wasmedge.ValType_I32
		}

		retType := []wasmedge.ValType{}
		if fn.Returns {
			retType = append(retType, wasmedge.ValType_I32)
		}
		funcType := wasmedge.NewFunctionType(argsType, retType)

		wasmEdgeHostFn := wasmedge.NewHostFunction(funcType, wasmHostFn, 0)
		imports.AddHostFunction(fn.Name, wasmEdgeHostFn)
	}
}
