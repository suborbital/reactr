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

//RetrySeconds returns an Option to set the worker retry seconds
func RetrySeconds(secs int) Option {
	return func(opts workerOpts) workerOpts {
		opts.retrySecs = secs
		return opts
	}
}

//MaxRetries returns an Option to set the worker maximum retry count
func MaxRetries(count int) Option {
	return func(opts workerOpts) workerOpts {
		opts.numRetries = count
		return opts
	}
}
