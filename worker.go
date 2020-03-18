package hive

const defaultChanSize = 1024

type worker struct {
	workChan chan Job
	runner   Runnable
}

// newWorker creates a new worker
func newWorker(runner Runnable) *worker {
	w := &worker{
		workChan: make(chan Job, defaultChanSize),
		runner:   runner,
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
