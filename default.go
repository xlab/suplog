package suplog

import (
	"context"
	"time"
)

var (
	// DefaultLogger provides a way to pass default logger into Hook initializers.
	//
	// nolint:gochecknoglobals
	DefaultLogger = &suplogger{}

	// Ensure that defaultOpt's type *suplogger matches Logger interface
	// during the compilation phase.
	_ Logger = DefaultLogger
)

// CLASSIC LOGGER METHODS

// Print will print to the underlying writer.
func Print(str string) {
	DefaultLogger.Print(str)
}

// Printf will print a formatted message to the underlying writer.
func Printf(format string, args ...interface{}) {
	DefaultLogger.Printf(format, args...)
}

// Notification will log a notification message.
func Notification(format string, args ...interface{}) {
	DefaultLogger.Notification(format, args...)
}

// Success will log a success message.
func Success(format string, args ...interface{}) {
	DefaultLogger.Success(format, args...)
}

// Warning will log a warning message.
func Warning(format string, args ...interface{}) {
	DefaultLogger.Warning(format, args...)
}

// Error will log an error message.
func Error(format string, args ...interface{}) {
	DefaultLogger.Error(format, args...)
}

// Debug will log a debug line.
func Debug(format string, args ...interface{}) {
	DefaultLogger.Debug(format, args...)
}

// OUTPUTTER METHODS
//
// Part A: Context providers

func WithField(key string, value interface{}) Logger {
	return DefaultLogger.WithField(key, value)
}

func WithFields(fields Fields) Logger {
	return DefaultLogger.WithFields(fields)
}

func WithError(err error) Logger {
	return DefaultLogger.WithError(err)
}

func WithContext(ctx context.Context) Logger {
	return DefaultLogger.WithContext(ctx)
}

func WithTime(t time.Time) Logger {
	return DefaultLogger.WithTime(t)
}

// Part B: Formatted logging methods

func Logf(level Level, format string, args ...interface{}) {
	DefaultLogger.Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) {
	DefaultLogger.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

func Warningf(format string, args ...interface{}) {
	DefaultLogger.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	DefaultLogger.Panicf(format, args...)
}

// Part C: Shortcut logging methods

func Log(level Level, args ...interface{}) {
	DefaultLogger.Log(level, args...)
}

func Trace(args ...interface{}) {
	DefaultLogger.Trace(args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	DefaultLogger.Panic(args...)
}

func Logln(level Level, args ...interface{}) {
	DefaultLogger.Logln(level, args...)
}

func Traceln(args ...interface{}) {
	DefaultLogger.Traceln(args...)
}

func Debugln(args ...interface{}) {
	DefaultLogger.Debugln(args...)
}

func Infoln(args ...interface{}) {
	DefaultLogger.Infoln(args...)
}

func Println(args ...interface{}) {
	DefaultLogger.Println(args...)
}

func Warningln(args ...interface{}) {
	DefaultLogger.Warningln(args...)
}

func Errorln(args ...interface{}) {
	DefaultLogger.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	DefaultLogger.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	DefaultLogger.Panicln(args...)
}

func FnName() string {
	return DefaultLogger.CallerName()
}
