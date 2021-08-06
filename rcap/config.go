package rcap

import "github.com/pkg/errors"

var ErrCapabilityNotEnabled = errors.New("capability is not enabled")

// CapabilityConfig is configuration for a Runnable's capabilities
type CapabilityConfig struct {
	LoggerConfig
	HTTPConfig
	CacheConfig
	FileConfig
	RequestHandlerConfig
}
