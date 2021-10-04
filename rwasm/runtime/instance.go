package runtime

import (
	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rt"
)

// WasmInstance is an instance of a Wasm runtime
type WasmInstance struct {
	runtime RuntimeInstance

	ctx *rt.Ctx

	ffiResult []byte

	resultChan chan []byte
	errChan    chan rt.RunErr
}

// RuntimeBuilder is a factory-style interface that can build Wasm runtimes
type RuntimeBuilder interface {
	New() (RuntimeInstance, error)
}

// RuntimeInstance is an interface that wraps various underlying Wasm runtimes like Wasmer, wasmTime
type RuntimeInstance interface {
	Call(fn string, args ...interface{}) (interface{}, error)
	ReadMemory(pointer int32, size int32) []byte
	WriteMemory(data []byte) (int32, error)
	WriteMemoryAtLocation(pointer int32, data []byte)
	Deallocate(pointer int32, length int)
	Close()
}

// instanceReference holds a reference to a particular WasmInstance
type instanceReference struct {
	Inst *WasmInstance
}

/////////////////////////////////////////////////////////////////////////////
// below is the wasm glue code used to manipulate wasm instance memory     //
// this requires a set of functions to be available within the wasm module //
// - allocate                                                              //
// - deallocate                                                            //
/////////////////////////////////////////////////////////////////////////////

// Call executes a function from the Wasm Module
func (w *WasmInstance) Call(fn string, args ...interface{}) (interface{}, error) {
	return w.runtime.Call(fn, args...)
}

// ExecutionResult gets the runnable's execution results
func (w *WasmInstance) ExecutionResult() ([]byte, error) {
	// determine if the instance called return_result or return_error
	select {
	case res := <-w.resultChan:
		return res, nil
	case err := <-w.errChan:
		return nil, err
	default:
		// do nothing and fall through
	}

	return nil, nil
}

// SendExecutionResult allows FFI functions to send the run result
func (w *WasmInstance) SendExecutionResult(result []byte, runErr *rt.RunErr) {
	if runErr != nil {
		w.errChan <- *runErr
	} else if result != nil {
		w.resultChan <- result
	}
}

// Ctx returns the instance's Ctx
func (w *WasmInstance) Ctx() *rt.Ctx {
	return w.ctx
}

func (w *WasmInstance) SetFFIResult(data []byte) error {
	if w.ffiResult != nil {
		return errors.New("instance ffiResult is already set")
	}

	w.ffiResult = data

	return nil
}

func (w *WasmInstance) UseFFIResult() ([]byte, error) {
	if w.ffiResult == nil {
		return nil, errors.New("instance ffiResult is not set")
	}

	defer func() {
		w.ffiResult = nil
	}()

	return w.ffiResult, nil
}

func (w *WasmInstance) ReadMemory(pointer int32, size int32) []byte {
	return w.runtime.ReadMemory(pointer, size)
}

func (w *WasmInstance) WriteMemory(data []byte) (int32, error) {
	return w.runtime.WriteMemory(data)
}

func (w *WasmInstance) WriteMemoryAtLocation(pointer int32, data []byte) {
	w.runtime.WriteMemoryAtLocation(pointer, data)
}

func (w *WasmInstance) Deallocate(pointer int32, length int) {
	w.runtime.Deallocate(pointer, length)
}
