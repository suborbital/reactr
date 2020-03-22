package hive

const defaultChanSize = 1024

type worker struct {
	workChan chan Job
	runner   Runnable
	options  workerOpts
}

type workerOpts struct {
	poolSize int
}

// newWorker creates a new worker
func newWorker(runner Runnable, opts workerOpts) *worker {
	w := &worker{
		workChan: make(chan Job, defaultChanSize),
		runner:   runner,
		options:  opts,
	}

	return w
}

func (w *worker) schedule(job Job) *Result {
	result := newResult()
	job.result = result

	go func() {
		w.workChan <- job
	}()

	return result
}

func (w *worker) start(runFunc RunFunc) {
	if w.workChan == nil {
		w.workChan = make(chan Job, defaultChanSize)
	}

	// fill the "pool" with goroutines
	for i := 0; i < w.options.poolSize; i++ {
		go func() {
			for {
				// wait for the next job
				job := <-w.workChan

				result, err := w.runner.Run(job, runFunc)
				if err != nil {
					job.result.sendErr(err)
				}

				job.result.sendResult(result)
			}
		}()
	}
}

func defaultOpts() workerOpts {
	o := workerOpts{
		poolSize: 1,
	}

	return o
}
