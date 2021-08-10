package rcap

import "github.com/pkg/errors"

var ErrCapabilityNotEnabled = errors.New("capability is not enabled")

// CapabilityConfig is configuration for a Runnable's capabilities
type CapabilityConfig struct {
	Logger         LoggerConfig         `json:"logger" yaml:"logger"`
	HTTP           HTTPConfig           `json:"http" yaml:"http"`
	GraphQL        GraphQLConfig        `json:"graphql" yaml:"graphql"`
	Auth           AuthConfig           `json:"auth" yaml:"auth"`
	Cache          CacheConfig          `json:"cache" yaml:"cache"`
	File           FileConfig           `json:"file" yaml:"file"`
	RequestHandler RequestHandlerConfig `json:"requestHandler" yaml:"requestHandler"`
}
