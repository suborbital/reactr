package hive

import (
	"encoding/json"

	"github.com/suborbital/hive/util"

	"github.com/pkg/errors"
)

// Result describes a result
type Result struct {
	ID         string
	resultChan chan interface{}
	errChan    chan error
}

// Discard returns immediately and discards the eventual results and thus prevents the memory from hanging around
func (r *Result) Discard() {
	go func() {
		select {
		case <-r.resultChan:
		case <-r.errChan:
		}
	}()
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
	res, err := r.Then()
	if err != nil {
		return 0, err
	}

	intVal, ok := res.(int)
	if !ok {
		return 0, errors.New("failed to convert result to Int")
	}

	return intVal, nil
}

// ThenJSON unmarshals the result or returns the error from a Result
func (r *Result) ThenJSON(out interface{}) error {
	res, err := r.Then()
	if err != nil {
		return err
	}

	b, ok := res.([]byte)
	if !ok {
		return errors.New("cannot unmarshal, result is not []byte")
	}

	if err := json.Unmarshal(b, out); err != nil {
		return errors.Wrap(err, "failed to Unmarshal result")
	}

	return nil
}

func newResult() *Result {
	r := &Result{
		ID:         util.GenerateResultID(),
		resultChan: make(chan interface{}, 1), // buffered, so the result can be written and related goroutines can end before Then() is called
		errChan:    make(chan error, 1),
	}

	return r
}

func (r *Result) sendResult(result interface{}) {
	// if the result is another Result,
	// wait for its result and recursively send it
	// or if the result is a group, wait on the
	// group and propogate the error if any
	if res, ok := result.(*Result); ok {
		go func() {
			if newResult, err := res.Then(); err != nil {
				r.sendErr(err)
			} else {
				r.sendResult(newResult)
			}
		}()

		return
	} else if grp, ok := result.(*Group); ok {
		go func() {
			if err := grp.Wait(); err != nil {
				r.sendErr(err)
			} else {
				r.sendResult(nil)
			}
		}()

		return
	}

	r.resultChan <- result
}

func (r *Result) sendErr(err error) {
	r.errChan <- err
}
