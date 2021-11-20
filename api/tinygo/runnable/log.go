//go:build tinygo.wasm

package runnable

// #include <reactr.h>
import "C"
import (
	"fmt"
)

func logAtLevel(message string, level LogLevel) {
	msgPtr, size := rawSlicePointer([]byte(message))

	C.log_msg(msgPtr, size, int32(level), ident())
}

type LogLevel int32

const (
	LogLevelError LogLevel = iota + 1
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

func Debug(message string) {
	logAtLevel(message, LogLevelDebug)
}

func Info(message string) {
	logAtLevel(message, LogLevelInfo)
}

func Warn(message string) {
	logAtLevel(message, LogLevelWarn)
}

func Error(message string) {
	logAtLevel(message, LogLevelError)
}

func Debugf(format string, args ...interface{}) {
	Debug(fmt.Sprintf(format, args...))
}

func Infof(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	Error(fmt.Sprintf(format, args...))
}
