package hive

import (
	"github.com/suborbital/gust/gapi"
)

//DoFunc is a function that runs a job of a predetermined type
type DoFunc func(interface{}) *Result

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

// Handle registers a Runnable with the Hive and returns a shortcut function to run those jobs
func (h *Hive) Handle(jobType string, runner Runnable, options ...Option) DoFunc {
	h.handle(jobType, runner, options...)

	helper := func(data interface{}) *Result {
		job := Job{
			jobType: jobType,
			data:    data,
		}

		return h.Do(job)
	}

	return helper
}

// Job is a shorter alias for NewJob
func (h *Hive) Job(jobType string, data interface{}) Job {
	return NewJob(jobType, data)
}

// Server returns a new Hive server
func (h *Hive) Server(opts ...gapi.OptionsModifier) *Server {
	return newServer(h, opts...)
}
