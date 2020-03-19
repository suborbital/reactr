package hive

import "errors"

// Result describes a result
type Result struct {
	resultChan chan interface{}
	errChan    chan error
}

// Then returns the result or error from a Result
func (r *Result) Then() (interface{}, error) {
	select {
	case res := <-r.resultChan:
		return res, nil
	case err := <-r.errChan:
		return nil, err
	}
}

// ThenInt returns the result or error from a Result
func (r *Result) ThenInt() (int, error) {
	select {
	case res := <-r.resultChan:
		intVal, ok := res.(int)
		if !ok {
			return 0, errors.New("failed to convert result to Int")
		}

		return intVal, nil
	case err := <-r.errChan:
		return 0, err
	}
}

func newResult() *Result {
	r := &Result{
		resultChan: make(chan interface{}),
		errChan:    make(chan error),
	}

	return r
}

func (r *Result) sendResult(result interface{}) {
	// if the result is another Result,
	// wait for its result and recursively send it
	if res, ok := result.(*Result); ok {
		go func() {
			if newResult, err := res.Then(); err != nil {
				r.sendErr(err)
			} else {
				r.sendResult(newResult)
			}
		}()

		return
	}

	r.resultChan <- result
}

func (r *Result) sendErr(err error) {
	r.errChan <- err
}
