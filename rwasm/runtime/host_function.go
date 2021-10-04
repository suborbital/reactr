package runtime

// HostFn describes a host function callable from within a Runnable module
type HostFn struct {
	Name     string
	ArgCount int
	Returns  bool
	HostFn   innerFunc
}

type innerFunc func(...interface{}) (interface{}, error)

// NewHostFn creates a new host function
func NewHostFn(name string, argCount int, returns bool, fn innerFunc) HostFn {
	h := HostFn{
		Name:     name,
		ArgCount: argCount,
		Returns:  returns,
		HostFn:   fn,
	}

	return h
}
