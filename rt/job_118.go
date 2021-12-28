//go:build go1.18

package rt

import (
	"github.com/google/uuid"
	"github.com/suborbital/reactr/request"
)

// Job describes a job to be done
type Job[T Input, R Output] struct {
	uuid    string
	jobType string
	result  *Result[R]
	data    T

	caps *Capabilities
	req  *request.CoordinatedRequest
}

// NewJob creates a new job
func NewJob[T Input, R Output](jobType string, data T) Job[T, R] {
	j := Job[T, R]{
		uuid:    uuid.New().String(),
		jobType: jobType,
		result: &Result[R]{},
		data:    data,
	}

	// TODO: figure this out for generics
	// // detect the coordinated request
	// if req, ok := data.(*request.CoordinatedRequest); ok {
	// 	j.req = req
	// }

	return j
}

// UUID returns the Job's UUID
func (j Job[T, R]) UUID() string {
	return j.uuid
}

// Unmarshal unmarshals the job's data into a struct
// func (j Job[T, R]) Unmarshal(target any) error {
// 	if bytes, ok := j.data.([]byte); ok {
// 		return json.Unmarshal(bytes, target)
// 	}

// 	return errors.New("failed to Unmarshal, job data is not []byte")
// }

// String returns the string value of a job's data
// func (j Job[T, R]) String() string {
// 	if s, isString := j.data.(string); isString {
// 		return s
// 	} else if b, isBytes := j.data.([]byte); isBytes {
// 		return string(b)
// 	}

// 	return ""
// }

// Bytes returns the []byte value of the job's data
// func (j Job[T, R]) Bytes() []byte {
// 	if v, ok := j.data.([]byte); ok {
// 		return v
// 	} else if s, ok := j.data.(string); ok {
// 		return []byte(s)
// 	}

// 	return nil
// }

// Int returns the int value of the job's data
// func (j Job[T, R]) Int() int {
// 	if v, ok := j.data.(int); ok {
// 		return v
// 	}

// 	return 0
// }

// Data returns the "raw" data for the job
func (j Job[T, R]) Data() T {
	return j.data
}

// Req returns the Coordinated request attached to the Job
func (j Job[T, R]) Req() *request.CoordinatedRequest {
	return j.req
}
