package hive

//RunFunc describes a function to schedule work
type RunFunc func(Job) *Result

// Runnable describes something that is runnable
type Runnable interface {
	Run(Job, RunFunc) (interface{}, error)
}
