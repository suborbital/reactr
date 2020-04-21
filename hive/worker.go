package hive

const defaultChanSize = 1024

// worker describes a worker
type worker interface {
	schedule(Job) *Result
	start(RunFunc) error
}

type goWorker struct {
	workChan chan Job
	runner   Runnable
	options  workerOpts
}

type workerOpts struct {
	poolSize int
}

// newGoWorker creates a new goWorker
func newGoWorker(runner Runnable, opts workerOpts) *goWorker {
	w := &goWorker{
		workChan: make(chan Job, defaultChanSize),
		runner:   runner,
		options:  opts,
	}

	return w
}

func (w *goWorker) schedule(job Job) *Result {
	result := newResult()
	job.result = result

	go func() {
		w.workChan <- job
	}()

	return result
}

func (w *goWorker) start(runFunc RunFunc) error {
	if w.workChan == nil {
		w.workChan = make(chan Job, defaultChanSize)
	}

	// fill the "pool" with goroutines
	for i := 0; i < w.options.poolSize; i++ {
		runnerCopy := w.runner

		go func() {
			for {
				// wait for the next job
				job := <-w.workChan

				result, err := runnerCopy.Run(job, runFunc)
				if err != nil {
					job.result.sendErr(err)
					continue
				}

				job.result.sendResult(result)
			}
		}()
	}

	return nil
}

func defaultOpts() workerOpts {
	o := workerOpts{
		poolSize: 1,
	}

	return o
}
