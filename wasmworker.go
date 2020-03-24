package hive

import (
	"github.com/pkg/errors"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

type wasmWorker struct {
	workChan chan Job
	runner   Runnable
	options  workerOpts
}

// newWasmWorker creates a new wasmWorker
func newWasmWorker(runner Runnable, opts workerOpts) *wasmWorker {
	w := &wasmWorker{
		workChan: make(chan Job, defaultChanSize),
		runner:   runner,
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

	// fill the "pool" with goroutines
	for i := 0; i < w.options.poolSize; i++ {
		runnerCopy := *w.runner.(*WasmRunner)

		instance, err := newInstance(runnerCopy.wasmFile)
		if err != nil {
			return errors.Wrap(err, "wasmWorker failed to newInstance")
		}

		runnerCopy.useInstance(instance)

		runner := &runnerCopy // make it a pointer again

		go func() {
			for {
				// wait for the next job
				job := <-w.workChan

				result, err := runner.Run(job, runFunc)
				if err != nil {
					job.result.sendErr(err)
				}

				job.result.sendResult(result)
			}
		}()
	}

	return nil
}

func newInstance(path string) (*wasm.Instance, error) {
	bytes, err := wasm.ReadBytes(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wasm.ReadBytes")
	}

	instance, err := wasm.NewInstance(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wasm.NewInstance")
	}

	return &instance, nil
}
