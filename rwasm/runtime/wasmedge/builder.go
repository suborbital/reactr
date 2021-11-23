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
	ast     *wasmedge.AST
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
	err := w.setupAST()
	if err != nil {
		return nil, err
	}

	// Create store
	store := wasmedge.NewStore()

	// Create executor
	executor := wasmedge.NewExecutor()

	// Register import object
	executor.RegisterImport(store, w.imports)

	wasiImports := wasmedge.NewWasiImportObject(nil, nil, nil, nil)
	executor.RegisterImport(store, wasiImports)

	// Instantiate store
	executor.Instantiate(store, w.ast)

	wasiStart := store.FindFunction("_start")
	if wasiStart != nil {
		if _, err := executor.Invoke(store, "_start"); err != nil {
			return nil, errors.Wrap(err, "failed to _start")
		}
	}
	init := store.FindFunction("init")
	if init != nil {
		if _, err := executor.Invoke(store, "init"); err != nil {
			return nil, errors.Wrap(err, "failed to init")
		}
	}

	inst := &WasmEdgeRuntime {
		store: store,
		executor: executor,
	}

	return inst, nil
}

func (w *WasmEdgeBuilder) setupAST() error {
	if w.ast == nil {
		// Set not to print debug info
		wasmedge.SetLogErrorLevel()

		moduleBytes, err := w.ref.Bytes()
		if err != nil {
			return errors.Wrap(err, "failed to get ref ModuleBytes")
		}

		// Create Loader
		loader := wasmedge.NewLoader()

		// Create AST
		ast, err := loader.LoadBuffer(moduleBytes)
		if err != nil {
			return errors.Wrap(err, "failed to create ast")
		}
		loader.Release()

		// Validate the ast
		val := wasmedge.NewValidator()
		err = val.Validate(ast)
		if err != nil {
			return errors.Wrap(err, "failed to validate ast")
		}
		defer val.Release()

		// Create import object
		imports := wasmedge.NewImportObject("env")

		// mount the Runnable API host functions to the module's imports
		addHostFns(imports, w.hostFns...)

		w.imports = imports
		w.ast = ast
	}

	return nil
}
