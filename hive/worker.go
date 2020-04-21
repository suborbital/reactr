package hive

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const defaultChanSize = 1024

// worker describes a worker
type worker interface {
	schedule(Job) *Result
	start(RunFunc) error
}

type goWorker struct {
	workChan chan Job
	runner   Runnable
	options  workerOpts
}

type workerOpts struct {
	poolSize   int
	numRetries int
	retrySecs  int
}

// newGoWorker creates a new goWorker
func newGoWorker(runner Runnable, opts workerOpts) *goWorker {
	w := &goWorker{
		workChan: make(chan Job, defaultChanSize),
		runner:   runner,
		options:  opts,
	}

	return w
}

func (w *goWorker) schedule(job Job) *Result {
	result := newResult()
	job.result = result

	go func() {
		w.workChan <- job
	}()

	return result
}

func (w *goWorker) start(runFunc RunFunc) error {
	if w.workChan == nil {
		w.workChan = make(chan Job, defaultChanSize)
	}

	started := 0
	attempts := 0

	for {
		// fill the "pool" with goroutines
		for i := started; i < w.options.poolSize; i++ {
			if err := w.runner.OnStart(); err != nil {
				fmt.Println(errors.Wrapf(err, "Runnable returned OnStart error, will retry in %ds", w.options.retrySecs))
				break
			} else {
				started++
			}

			go func() {
				for {
					// wait for the next job
					job := <-w.workChan

					result, err := w.runner.Run(job, runFunc)
					if err != nil {
						job.result.sendErr(err)
						continue
					}

					job.result.sendResult(result)
				}
			}()
		}

		if started == w.options.poolSize {
			break
		} else {
			if attempts >= w.options.numRetries {
				return fmt.Errorf("attempted to start worker %d times, Runnable returned error each time", w.options.numRetries)
			}

			attempts++
			<-time.After(time.Duration(time.Second * time.Duration(w.options.retrySecs)))
		}
	}

	return nil
}

func defaultOpts() workerOpts {
	o := workerOpts{
		poolSize:   1,
		retrySecs:  3,
		numRetries: 5,
	}

	return o
}
