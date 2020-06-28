package hive

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

type scheduler struct {
	workers map[string]*worker

	starter sync.Once
	sync.Mutex
}

func newScheduler() *scheduler {
	s := &scheduler{
		workers: map[string]*worker{},
		Mutex:   sync.Mutex{},
	}

	return s
}

func (s *scheduler) schedule(job Job) *Result {
	s.starter.Do(func() {
		if s.workers == nil {
			s.workers = map[string]*worker{}
		}
	})

	result := newResult()

	worker := s.getWorker(job.jobType)
	if worker == nil {
		result.sendErr(fmt.Errorf("failed to getRunnable for jobType %q", job.jobType))
		return result
	}

	go func() {
		if !worker.isStarted() {
			// "recursively" pass this function as the runFunc for the runnable
			if err := worker.start(s.schedule); err != nil {
				result.sendErr(errors.Wrapf(err, "failed start worker for jobType %q", job.jobType))
				return
			}
		}

		job.result = result
		worker.schedule(job)
	}()

	return result
}

// handle adds a handler
func (s *scheduler) handle(jobType string, runnable Runnable, options ...Option) {
	s.Lock()
	defer s.Unlock()

	// apply the provided options
	opts := defaultOpts(jobType)
	for _, o := range options {
		opts = o(opts)
	}

	w := newWorker(runnable, opts)
	if s.workers == nil {
		s.workers = map[string]*worker{jobType: w}
	} else {
		s.workers[jobType] = w
	}
}

func (s *scheduler) getWorker(jobType string) *worker {
	s.Lock()
	defer s.Unlock()

	if s.workers == nil {
		return nil
	}

	if w, ok := s.workers[jobType]; ok {
		return w
	}

	return nil
}
