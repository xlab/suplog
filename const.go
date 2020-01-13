package suplog

import (
	"fmt"
	"strings"
)

// These are the different logging levels.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// ParseLevel takes a string level and returns the output log level constant.
func ParseLevel(levelName string) (level Level, err error) {
	switch strings.ToLower(levelName) {
	case "panic":
		level = PanicLevel
	case "fatal":
		level = FatalLevel
	case "error":
		level = ErrorLevel
	case "warn", "warning":
		level = WarnLevel
	case "info":
		level = InfoLevel
	case "debug":
		level = DebugLevel
	case "trace":
		level = TraceLevel
	default:
		err = fmt.Errorf("not a valid output Level: %s", levelName)
	}

	return
}
