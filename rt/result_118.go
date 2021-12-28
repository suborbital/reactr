//go:build go1.18

package rt

import (
)

// Result describes a result
type Result[R any] struct {
	uuid string
	data R
	err  error

	resultChan chan bool
	errChan    chan bool
}

// ResultFunc is a result callback function.
type ResultFunc[R any] func(R, error)

func newResult[R any](uuid string) *Result[R] {
	r := &Result[R]{
		uuid:       uuid,
		resultChan: make(chan bool, 1), // buffered, so the result can be written and related goroutines can end before Then() is called
		errChan:    make(chan bool, 1),
	}

	return r
}

// UUID returns the result/job's UUID
func (r *Result[R]) UUID() string {
	return r.uuid
}

// Then returns the result or error from a Result
func (r *Result[R]) Then() (R, error) {
	select {
	case <-r.resultChan:
		return r.data, nil
	case <-r.errChan:
		// TODO: determine how to return nil instead of the should-be-nil r.data here
		return r.data, r.err
	}
}

// ThenInt returns the result or error from a Result
// func (r *Result[R]) ThenInt() (int, error) {
// 	res, err := r.Then()
// 	if err != nil {
// 		return 0, err
// 	}

// 	intVal, ok := res.(int)
// 	if !ok {
// 		return 0, errors.New("failed to convert result to Int")
// 	}

// 	return intVal, nil
// }

// ThenJSON unmarshals the result or returns the error from a Result
// func (r *Result[R]) ThenJSON(out any) error {
// 	res, err := r.Then()
// 	if err != nil {
// 		return err
// 	}

// 	b, ok := res.([]byte)
// 	if !ok {
// 		return errors.New("cannot unmarshal, result is not []byte")
// 	}

// 	if err := json.Unmarshal(b, out); err != nil {
// 		return errors.Wrap(err, "failed to Unmarshal result")
// 	}

// 	return nil
// }

// ThenDo accepts a callback function to be called asynchronously when the result completes.
func (r *Result[R]) ThenDo(do ResultFunc[R]) {
	go func() {
		res, err := r.Then()
		do(res, err)
	}()
}

// Discard returns immediately and discards the eventual results and thus prevents the memory from hanging around
func (r *Result[R]) Discard() {
	go func() {
		r.Then()
	}()
}

func (r *Result[R]) sendResult(data R) {
	// if the result is another Result,
	// wait for its result and recursively send it
	// or if the result is a group, wait on the
	// group and propogate the error if any
	// if res, ok := data.(*Result[R]); ok {
	// 	go func() {
	// 		if newResult, err := res.Then(); err != nil {
	// 			r.sendErr(err)
	// 		} else {
	// 			r.sendResult(newResult)
	// 		}
	// 	}()

	// 	return
	// } 
	// else if grp, ok := data.(*Group[T, R]); ok {
	// 	go func() {
	// 		if err := grp.Wait(); err != nil {
	// 			r.sendErr(err)
	// 		} else {
	// 			r.sendResult(nil)
	// 		}
	// 	}()

	// 	return
	// }

	r.data = data
	r.resultChan <- true
}

func (r *Result[R]) sendErr(err error) {
	r.err = err
	r.errChan <- true
}
