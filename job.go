package hive

import (
	"encoding/json"
	"errors"
)

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

// Unmarshal unmarshals the job's data into a struct
func (j *Job) Unmarshal(target interface{}) error {
	if bytes, ok := j.data.([]byte); ok {
		return json.Unmarshal(bytes, target)
	}

	return errors.New("failed to Unmarshal, job data is not []byte")
}

func (j *Job) String() string {
	if s, ok := j.data.(string); ok {
		return s
	}

	return ""
}

// Bytes returns the []byte value of the job's data
func (j *Job) Bytes() []byte {
	if v, ok := j.data.([]byte); ok {
		return v
	}

	return nil
}

// Int returns the int value of the job's data
func (j *Job) Int() int {
	if v, ok := j.data.(int); ok {
		return v
	}

	return 0
}

// Data returns the "raw" data for the job
func (j *Job) Data() interface{} {
	return j.data
}
