package rt

import (
	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rcap"
)

var ErrCapabilityNotAvailable = errors.New("capability not available")

// Capabilities define the capabilities available to a Runnable
type Capabilities struct {
	Cache rcap.Cache

	doFunc coreDoFunc
}
