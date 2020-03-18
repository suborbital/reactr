package hive

// Hive represents the main control object
type Hive struct {
	*scheduler
}

// New returns a Hive ready to accept Jobs
func New() *Hive {
	h := &Hive{
		scheduler: newScheduler(),
	}

	return h
}

// Do schedules a job to be worked on and returns a result object
func (h *Hive) Do(job Job) *Result {
	return h.schedule(job)
}
