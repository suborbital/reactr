package runtimewasmedge

import (
	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rwasm/moduleref"
	"github.com/suborbital/reactr/rwasm/runtime"
	"github.com/second-state/WasmEdge-go/wasmedge"
)

// WasmEdgeBuilder is a WasmEdge implementation of the instanceBuilder interface
type WasmEdgeBuilder struct {
	ref     *moduleref.WasmModuleRef
	hostFns []runtime.HostFn
	imports *wasmedge.ImportObject
}

// NewBuilder create a new WasmEdgeBuilder
func NewBuilder(ref *moduleref.WasmModuleRef, hostFns ...runtime.HostFn) runtime.RuntimeBuilder {
	w := &WasmEdgeBuilder {
		ref:     ref,
		hostFns: hostFns,
	}
	return w
}

func (w *WasmEdgeBuilder) New() (runtime.RuntimeInstance, error) {
	imports, _ := w.internals()

	moduleBytes, err := w.ref.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ref ModuleBytes")
	}

	// Create configure
	conf := wasmedge.NewConfigure(wasmedge.WASI)

	// Create store
	store := wasmedge.NewStore()
	
	// Create VM by configure and external store
	vm := wasmedge.NewVMWithConfigAndStore(conf, store)

	// Register import object
	vm.RegisterImport(imports)

	// Instantiate wasm
	vm.LoadWasmBuffer(moduleBytes)
	vm.Validate()
	vm.Instantiate()

	wasiStart := store.FindFunction("_start")
	if wasiStart != nil {
		if _, err := vm.Execute("_start"); err != nil {
			return nil, errors.Wrap(err, "failed to _start")
		}
	}
	init := store.FindFunction("init")
	if init != nil {
		if _, err := vm.Execute("init"); err != nil {
			return nil, errors.Wrap(err, "failed to init")
		}
	}

	inst := &WasmEdgeRuntime {
		store: store,
		vm: vm,
	}

	return inst, nil
}

func (w *WasmEdgeBuilder) internals() (*wasmedge.ImportObject, error) {
	if w.imports == nil {
		// Set not to print debug info
		wasmedge.SetLogErrorLevel()

		// Create import object
		imports := wasmedge.NewImportObject("env")

		// mount the Runnable API host functions to the module's imports
		addHostFns(imports, w.hostFns...)

		w.imports = imports
	}

	return w.imports, nil
}
