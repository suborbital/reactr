package hive

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultChanSize = 256
)

// ErrJobTimeout and others are errors related to workers
var (
	ErrJobTimeout = errors.New("job timeout")
)

type worker struct {
	runner   Runnable
	workChan chan Job
	options  workerOpts

	threads    []*workThread
	threadLock sync.Mutex

	started bool
	starter sync.Once
}

// newWorker creates a new goWorker
func newWorker(runner Runnable, opts workerOpts) *worker {
	w := &worker{
		runner:     runner,
		workChan:   make(chan Job, defaultChanSize),
		options:    opts,
		threads:    make([]*workThread, opts.poolSize),
		threadLock: sync.Mutex{},
		started:    false,
	}

	return w
}

func (w *worker) schedule(job Job) {
	go func() {
		w.workChan <- job
	}()
}

func (w *worker) start(runFunc RunFunc) error {
	w.starter.Do(func() { w.started = true })

	started := 0
	attempts := 0

	for {
		// fill the "pool" with workThreads
		for i := started; i < w.options.poolSize; i++ {
			wt := newWorkThread(w.runner, w.workChan, w.options.jobTimeoutSeconds)

			// give the runner opportunity to provision resources if needed
			if err := w.runner.OnStart(); err != nil {
				fmt.Println(errors.Wrapf(err, "Runnable returned OnStart error, will retry in %ds", w.options.retrySecs))
				break
			} else {
				started++
			}

			wt.run(runFunc)

			w.threads[i] = wt
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

type workThread struct {
	runner         Runnable
	workChan       chan Job
	timeoutSeconds int
	ctx            context.Context
	cancelFunc     context.CancelFunc
}

func newWorkThread(runner Runnable, workChan chan Job, timeoutSeconds int) *workThread {
	ctx, cancelFunc := context.WithCancel(context.Background())

	wt := &workThread{
		runner:         runner,
		workChan:       workChan,
		timeoutSeconds: timeoutSeconds,
		ctx:            ctx,
		cancelFunc:     cancelFunc,
	}

	return wt
}

func (wt *workThread) run(runFunc RunFunc) {
	go func() {
		for {
			// die if the context has been cancelled
			if wt.ctx.Err() != nil {
				break
			}

			// wait for the next job
			job := <-wt.workChan

			var result interface{}
			var err error

			if wt.timeoutSeconds == 0 {
				result, err = wt.runner.Run(job, runFunc)
			} else {
				result, err = wt.runWithTimeout(job, runFunc)
			}

			if err != nil {
				job.result.sendErr(err)
				continue
			}

			job.result.sendResult(result)
		}
	}()
}

func (wt *workThread) runWithTimeout(job Job, runFunc RunFunc) (interface{}, error) {
	resultChan := make(chan interface{})
	errChan := make(chan error)

	go func() {
		result, err := wt.runner.Run(job, runFunc)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- result
		}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(time.Duration(time.Second * time.Duration(wt.timeoutSeconds))):
		return nil, ErrJobTimeout
	}
}

func (wt *workThread) Stop() {
	wt.cancelFunc()
}

type workerOpts struct {
	jobType           string
	poolSize          int
	jobTimeoutSeconds int
	numRetries        int
	retrySecs         int
}

func defaultOpts(jobType string) workerOpts {
	o := workerOpts{
		jobType:           jobType,
		poolSize:          1,
		jobTimeoutSeconds: 0,
		retrySecs:         3,
		numRetries:        5,
	}

	return o
}
