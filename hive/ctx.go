package hive

import "github.com/pkg/errors"

var errDoFuncNotSet = errors.New("do func has not been set")

// Ctx is a Job context
type Ctx struct {
	doFunc DoFunc
}

// Do runs a new job
func (c *Ctx) Do(job Job) *Result {
	if c.doFunc == nil {
		r := newResult(job.uuid, func(_ string) {})
		r.sendErr(errDoFuncNotSet)
		return r
	}

	return c.doFunc(job)
}
