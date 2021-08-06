package rcap

import "github.com/suborbital/vektor/vlog"

// LoggerConfig is configuration for the logger capability
type LoggerConfig struct {
	Enabled bool
}

// LoggerCapability provides a logger to Runnables
type LoggerCapability interface {
	Logger() *vlog.Logger
}

type loggerSource struct {
	config LoggerConfig
	log    *vlog.Logger
}

// DefaultLoggerSource returns a LoggerSource that provides vlog.Default
func DefaultLoggerSource(config LoggerConfig, logger *vlog.Logger) LoggerCapability {
	l := &loggerSource{
		config: config,
		log:    logger,
	}

	return l
}

// Logger returns the logger
func (l *loggerSource) Logger() *vlog.Logger {
	if !l.config.Enabled {
		return nil
	}

	return l.log
}
