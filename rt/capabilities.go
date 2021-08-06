package rt

import (
	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rcap"
	"github.com/suborbital/vektor/vlog"
)

var ErrCapabilityNotAvailable = errors.New("capability not available")

// Capabilities define the capabilities available to a Runnable
type Capabilities struct {
	Auth          rcap.AuthCapability
	LoggerSource  rcap.LoggerCapability
	HTTPClient    rcap.HTTPCapability
	GraphQLClient rcap.GraphQLCapability
	FileSource    rcap.FileCapability
	Cache         rcap.CacheCapability

	// RequestHandler and doFunc are special because they are more
	// sensitive; they could cause memory leaks or expose internal state,
	// so they cannot be swapped out for a different implementation.
	RequestHandler rcap.RequestHandlerCapability
	doFunc         coreDoFunc
}

func defaultCaps(logger *vlog.Logger) Capabilities {
	caps := Capabilities{
		Auth:          rcap.DefaultAuthProvider(rcap.AuthConfig{Enabled: true, Headers: nil}), // no authentication config is set up by default
		LoggerSource:  rcap.DefaultLoggerSource(rcap.LoggerConfig{Enabled: true}, logger),
		HTTPClient:    rcap.DefaultHTTPClient(rcap.HTTPConfig{Enabled: true}),
		GraphQLClient: rcap.DefaultGraphQLClient(rcap.GraphQLConfig{Enabled: true}),
		FileSource:    rcap.DefaultFileSource(rcap.FileConfig{Enabled: true}, nil), // set file access to nil by default, it can be set later.
		Cache:         rcap.DefaultCache(rcap.CacheConfig{Enabled: true}),

		// RequestHandler and doFunc don't get set here since they are set by
		// the rt and rwasm internals; a better solution for this should probably be found
	}

	return caps
}
