package hive

//RunFunc describes a function to schedule work
type RunFunc func(Job) *Result

// Runnable describes something that is runnable
type Runnable interface {
	// Run is the entrypoint for jobs handled by a Runnable
	Run(Job, RunFunc) (interface{}, error)

	// OnStart is called by the scheduler when a worker is started that will use the Runnable
	// OnStart will be called once for each worker in a pool
	OnStart() error
}
