package hive

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const defaultChanSize = 1024

type worker struct {
	workChan chan Job
	runner   Runnable
	options  workerOpts

	started bool
	starter sync.Once
}

// newWorker creates a new goWorker
func newWorker(runner Runnable, opts workerOpts) *worker {
	w := &worker{
		workChan: make(chan Job, defaultChanSize),
		runner:   runner,
		options:  opts,
		started:  false,
	}

	return w
}

func (w *worker) schedule(job Job) {
	go func() {
		w.workChan <- job
	}()
}

func (w *worker) start(runFunc RunFunc) error {
	w.starter.Do(func() {
		w.started = true

		if w.workChan == nil {
			w.workChan = make(chan Job, defaultChanSize)
		}
	})

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

func (w *worker) isStarted() bool {
	return w.started
}

type workerOpts struct {
	jobType    string
	poolSize   int
	numRetries int
	retrySecs  int
}

func defaultOpts(jobType string) workerOpts {
	o := workerOpts{
		jobType:    jobType,
		poolSize:   1,
		retrySecs:  3,
		numRetries: 5,
	}

	return o
}
