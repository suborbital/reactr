package rt

import "github.com/pkg/errors"

var (
	ErrSuperfuncMissing = errors.New("superfunc is missing")
)

// Superfunc is an abstraction that allows Reactr to run statically typed functions
type Superfunc func(*Ctx) (interface{}, error)

// superfuncRunner is a shim that allows a superfunc to be run as a Runnable
type superfuncRunner struct{}

func (s *superfuncRunner) Run(job Job, ctx *Ctx) (interface{}, error) {
	if job.sf == nil {
		return nil, ErrSuperfuncMissing
	}

	return job.sf(ctx)
}

func (s *superfuncRunner) OnChange(_ ChangeEvent) error { return nil }
