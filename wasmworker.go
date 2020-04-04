package hive

import (
	"github.com/pkg/errors"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

type wasmWorker struct {
	workChan chan Job
	runnable Runnable
	options  workerOpts
}

// newWasmWorker creates a new wasmWorker
func newWasmWorker(runnable Runnable, opts workerOpts) *wasmWorker {
	w := &wasmWorker{
		workChan: make(chan Job, defaultChanSize),
		runnable: runnable,
		options:  opts,
	}

	return w
}

func (w *wasmWorker) schedule(job Job) *Result {
	result := newResult()
	job.result = result

	go func() {
		w.workChan <- job
	}()

	return result
}

func (w *wasmWorker) start(runFunc RunFunc) error {
	if w.workChan == nil {
		w.workChan = make(chan Job, defaultChanSize)
	}

	wasmBytes, err := w.runnable.(*WasmRunner).WasmBytes()
	if err != nil {
		return errors.Wrap(err, "failed to WasmBytes")
	}

	// fill the "pool" with goroutines
	for i := 0; i < w.options.poolSize; i++ {
		runnableCopy := *w.runnable.(*WasmRunner)

		instance, err := wasm.NewInstance(wasmBytes)
		if err != nil {
			return errors.Wrap(err, "failed to wasm.NewInstance")
		}

		runnableCopy.useInstance(&instance)

		runnable := &runnableCopy // make it a pointer again

		go func() {
			for {
				// wait for the next job
				job := <-w.workChan

				result, err := runnable.Run(job, runFunc)
				if err != nil {
					job.result.sendErr(err)
				}

				job.result.sendResult(result)
			}
		}()
	}

	return nil
}
