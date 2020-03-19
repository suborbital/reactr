package hive

import (
	"fmt"
	"sync"
)

type scheduler struct {
	handler
	workers map[string]*worker
	sync.Mutex
}

func newScheduler() *scheduler {
	s := &scheduler{
		handler: handler{
			registered: map[string]Runnable{},
		},
		workers: map[string]*worker{},
		Mutex:   sync.Mutex{},
	}

	return s
}

func (s *scheduler) schedule(job Job) *Result {
	if s.workers == nil {
		s.workers = map[string]*worker{}
	}

	w, ok := s.workers[job.jobType]
	if !ok {
		runner := s.handler.getRunnable(job.jobType)
		if runner == nil {
			result := newResult()
			result.sendErr(fmt.Errorf("failed to getRunnable for jobType %q", job.jobType))
			return result
		}

		newWorker := newWorker(runner)
		newWorker.start(s.schedule) // "recursively" pass this function as the runFunc for the runnable

		s.Lock()
		s.workers[job.jobType] = newWorker
		s.Unlock()

		w = newWorker
	}

	return w.schedule(job)
}

type handler struct {
	registered map[string]Runnable
	sync.Mutex
}

// handle adds a handler
func (h *handler) handle(jobType string, runnable Runnable) {
	h.Lock()
	defer h.Unlock()

	if h.registered == nil {
		h.registered = map[string]Runnable{jobType: runnable}
	} else {
		h.registered[jobType] = runnable
	}
}

func (h *handler) getRunnable(jobType string) Runnable {
	h.Lock()
	defer h.Unlock()

	if h.registered == nil {
		return nil
	}

	if r, ok := h.registered[jobType]; ok {
		return r
	}

	return nil
}
