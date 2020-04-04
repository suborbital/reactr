package hive

import (
	"strings"

	"github.com/suborbital/hivew/hivew/util"

	"github.com/pkg/errors"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

//WasmRunner represents a wasm-based runnable
type WasmRunner struct {
	wasmFile string
	raw      *util.RawWASM
	inst     *wasm.Instance
}

// NewWasm returns a new WasmRunner
func NewWasm(path string) *WasmRunner {
	w := &WasmRunner{
		wasmFile: path,
	}

	return w
}

func newWasmFromRaw(raw *util.RawWASM) *WasmRunner {
	w := &WasmRunner{
		raw: raw,
	}

	return w
}

// Run runs a wasmRunner
func (w *WasmRunner) Run(job Job, run RunFunc) (interface{}, error) {
	if w.inst == nil {
		return nil, errors.New("WasmRunner attempted to Run with nil Instance")
	}

	input, ok := job.Data().(string)
	if !ok {
		return nil, errors.New("failed to run WASM job, input is not string")
	}

	inPointer := writeInput(w.inst, input)

	wasmRun := w.inst.Exports["run_e"]

	res, err := wasmRun(inPointer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wasmRun")
	}

	output := readOutput(w.inst, res.ToI32())

	// deallocate the memory used for the input and output
	deallocate(w.inst, inPointer, len(input))
	deallocate(w.inst, res.ToI32(), len(output))

	return output, nil
}

// WasmBytes returns the raw bytes of the runner's wasm
func (w *WasmRunner) WasmBytes() ([]byte, error) {
	if w.raw != nil {
		return w.raw.Contents, nil
	}

	return wasm.ReadBytes(w.wasmFile)
}

func (w *WasmRunner) useInstance(i *wasm.Instance) {
	w.inst = i
}

func writeInput(inst *wasm.Instance, input string) int32 {
	lengthOfInput := len(input)

	// Allocate memory for the input, and get a pointer to it.
	allocateResult, _ := inst.Exports["allocate_input"](lengthOfInput)
	inputPointer := allocateResult.ToI32()

	// Write the input into the memory.
	memory := inst.Memory.Data()[inputPointer:]

	for nth := 0; nth < lengthOfInput; nth++ {
		memory[nth] = input[nth]
	}

	// C-string terminates by NULL.
	memory[lengthOfInput] = 0

	return inputPointer
}

func readOutput(inst *wasm.Instance, pointer int32) string {
	memory := inst.Memory.Data()[pointer:]

	nth := 0
	var output strings.Builder

	for {
		if memory[nth] == 0 {
			break
		}

		output.WriteByte(memory[nth])
		nth++
	}

	return output.String()
}

func deallocate(inst *wasm.Instance, pointer int32, length int) {
	dealloc := inst.Exports["deallocate"]

	dealloc(pointer, length)
}
