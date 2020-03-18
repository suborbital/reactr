package hive

// Job describes a job to be done
type Job struct {
	jobType string
	data    interface{}
	result  *Result
}

// NewJob creates a new job
func NewJob(jobType string, data interface{}) Job {
	j := Job{
		jobType: jobType,
		data:    data,
	}

	return j
}

func (j *Job) String() string {
	if s, ok := j.data.(string); ok {
		return s
	}

	return ""
}
