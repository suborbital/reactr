package rt

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/suborbital/vektor/vlog"
)

type scheduler struct {
	workers map[string]*worker
	watcher *watcher
	cache   Cache
	logger  *vlog.Logger
	lock    sync.Mutex
}

func newScheduler(logger *vlog.Logger, cache Cache) *scheduler {
	s := &scheduler{
		workers: map[string]*worker{},
		cache:   cache,
		logger:  logger,
		lock:    sync.Mutex{},
	}

	s.watcher = newWatcher(s.schedule)

	return s
}

func (s *scheduler) schedule(job Job) *Result {
	result := newResult(job.UUID())

	worker := s.getWorker(job.jobType)
	if worker == nil {
		// if the job contains a task, let's auto-handle it
		if job.task != nil {
			s.handle(job.jobType, &taskRunner{})
			worker = s.getWorker(job.jobType)
		} else {
			result.sendErr(fmt.Errorf("failed to getWorker for jobType %q", job.jobType))
			return result
		}
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
	s.lock.Lock()
	defer s.lock.Unlock()

	// apply the provided options
	opts := defaultOpts(jobType)
	for _, o := range options {
		opts = o(opts)
	}

	w := newWorker(runnable, s.cache, opts)

	s.workers[jobType] = w

	if opts.preWarm {
		go func() {
			if err := w.start(s.schedule); err != nil {
				s.logger.Error(errors.Wrapf(err, "failed to preWarm %s worker", jobType))
			}
		}()
	}
}

func (s *scheduler) watch(sched Schedule) {
	s.watcher.watch(sched)
}

func (s *scheduler) getWorker(jobType string) *worker {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.workers == nil {
		return nil
	}

	if w, ok := s.workers[jobType]; ok {
		return w
	}

	return nil
}
