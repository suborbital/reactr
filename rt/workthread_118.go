//go:build go1.18

package rt

import (
	"context"
	"time"
)

type workThread[T any, R any] struct {
	runner         Runnable[T, R]
	workChan       chan *Job[T, R]
	timeoutSeconds int
	context        context.Context
	cancelFunc     context.CancelFunc
}

func newWorkThread[T any, R any](runner Runnable[T, R], workChan chan *Job[T, R], timeoutSeconds int) *workThread[T, R] {
	ctx, cancelFunc := context.WithCancel(context.Background())

	wt := &workThread[T, R]{
		runner:         runner,
		workChan:       workChan,
		timeoutSeconds: timeoutSeconds,
		context:        ctx,
		cancelFunc:     cancelFunc,
	}

	return wt
}

func (wt *workThread[T, R]) run() {
	go func() {
		for {
			// die if the context has been cancelled
			if wt.context.Err() != nil {
				break
			}

			// wait for the next job
			job := <-wt.workChan
			var err error

			ctx := newCtx(job.caps)

			var result R

			if wt.timeoutSeconds == 0 {
				// we pass in a dereferenced job so that the Runner cannot modify it
				result, err = wt.runner.Run(*job, ctx)
			} else {
				result, err = wt.runWithTimeout(job, ctx)
			}

			if err != nil {
				job.result.sendErr(err)
				continue
			}

			job.result.sendResult(result)
		}
	}()
}

func (wt *workThread[T, R]) runWithTimeout(job *Job[T, R], ctx *Ctx) (R, error) {
	resultChan := make(chan R)
	errChan := make(chan error)

	go func() {
		// we pass in a dereferenced job so that the Runner cannot modify it
		result, err := wt.runner.Run(*job, ctx)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- result
		}
	}()

	var empty R

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		// TODO: figure out how to return nil instead of result
		return empty, err
	case <-time.After(time.Duration(time.Second * time.Duration(wt.timeoutSeconds))):
		return empty, ErrJobTimeout
	}
}

func (wt *workThread[T, R]) Stop() {
	wt.cancelFunc()
}
