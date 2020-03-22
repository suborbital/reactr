package hive

// Option is a function that modifies workerOpts
type Option func(workerOpts) workerOpts

//PoolSize returns an Option to set the worker pool size
func PoolSize(size int) Option {
	return func(opts workerOpts) workerOpts {
		opts.poolSize = size
		return opts
	}
}
