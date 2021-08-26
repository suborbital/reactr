package rt

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/suborbital/vektor/vlog"
)

// coreDoFunc is an internal version of DoFunc that takes a
// Job pointer instead of a Job value for the best memory usage
type coreDoFunc func(job *Job) *Result

// core is the 'core scheduler' for reactr, handling execution of
// Tasks, Jobs, and Schedules
type core struct {
	workers map[string]*worker
	watcher *watcher
	log     *vlog.Logger
	lock    sync.RWMutex
}

func newCore(log *vlog.Logger) *core {
	c := &core{
		workers: map[string]*worker{},
		log:     log,
		lock:    sync.RWMutex{},
	}

	c.watcher = newWatcher(c.do)

	return c
}

func (c *core) do(job *Job) *Result {
	result := newResult(job.UUID())

	worker := c.findWorker(job.jobType)
	if worker == nil {
		result.sendErr(fmt.Errorf("failed to getWorker for jobType %q", job.jobType))
		return result
	}

	go func() {
		job.result = result

		worker.schedule(job)
	}()

	return result
}

// register adds a handler
func (c *core) register(jobType string, runnable Runnable, caps Capabilities, options ...Option) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// apply the provided options
	opts := defaultOpts(jobType)
	for _, o := range options {
		opts = o(opts)
	}

	w := newWorker(runnable, caps, opts)

	c.workers[jobType] = w

	go func() {
		if err := w.start(); err != nil {
			c.log.Error(errors.Wrapf(err, "failed to start %s worker", jobType))
		}
	}()
}

func (c *core) deRegister(jobType string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	worker, exists := c.workers[jobType]
	if !exists {
		// make this a no-op
		return nil
	}

	delete(c.workers, jobType)

	if err := worker.stop(); err != nil {
		return errors.Wrap(err, "failed to worker.stop")
	}

	return nil
}

func (c *core) watch(sched Schedule) {
	c.watcher.watch(sched)
}

func (c *core) findWorker(jobType string) *worker {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.workers == nil {
		return nil
	}

	if w, ok := c.workers[jobType]; ok {
		return w
	}

	return nil
}

func (c *core) hasWorker(jobType string) bool {
	w := c.findWorker(jobType)

	return w != nil
}
