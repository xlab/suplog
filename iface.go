package suplog

import (
	"context"
	"io"
	"time"
)

// Logger represents a full suplogger interface.
// It was inspired by previous ClassicLogger interface that we have to support,
// also logrus capabilities that are added here just recently.
type Logger interface {
	// Core logging methods

	Success(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})

	// Logrus context providers

	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger
	WithContext(ctx context.Context) Logger
	WithTime(t time.Time) Logger

	// Logrus formatted logging methods

	Logf(level Level, format string, args ...interface{})
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	// Logrus shortcut logging methods

	Log(level Level, args ...interface{})
	Trace(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Logln(level Level, args ...interface{})
	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

type LoggerConfigurator interface {
	SetFormatter(formatter Formatter)
	SetOutput(suplog io.Writer)
	SetLevel(level Level)
	GetLevel() Level
	IsLevelEnabled(level Level) bool
	AddHook(hook Hook)
	ReplaceHooks(hooks LevelHooks) LevelHooks
	SetStackTraceOffset(offset int)
	CallerName() string
}

var (
	_ StdLogger = &suplogger{}
	_ StdLogger = &Entry{}
)

// StdLogger is what your suplog-enabled library should take, that way
// it'll accept a stdlib logger (*log.Logger) and an suplog.Logger. There's no standard
// interface, this is the closest we get, unfortunately.
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}
