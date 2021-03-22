package suplog

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	blobHook "github.com/xlab/suplog/hooks/blob"
	bugsnagHook "github.com/xlab/suplog/hooks/bugsnag"
	debugHook "github.com/xlab/suplog/hooks/debug"

	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"github.com/xlab/suplog/stackcache"
)

// NewLogger constructs a new suplogger.
func NewLogger(wr io.Writer, formatter Formatter, hooks ...Hook) Logger {
	if formatter == nil {
		formatter = new(TextFormatter)
	}

	log := &suplogger{
		logger: &logrus.Logger{
			Out:       wr,
			Formatter: formatter,
			Hooks:     make(LevelHooks),
			Level:     DebugLevel,
			ExitFunc:  closer.Exit,
		},

		writer:           wr,
		mux:              new(sync.Mutex),
		stackTraceOffset: 0,
		initDone:         true,
	}

	log.reloadStackTraceCache()
	log.entry = log.logger.WithContext(context.Background())

	for _, h := range hooks {
		log.AddHook(h)
	}

	return log
}

type suplogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry

	mux              *sync.Mutex
	writer           io.Writer
	stack            stackcache.StackCache
	stackTraceOffset int

	init     sync.Once
	initDone bool
	closed   bool
}

func (l *suplogger) initOnce() {
	l.init.Do(func() {
		if l.initDone {
			// bail out if init already done (if New contstructor has been used).
			return
		}
		if l.writer == nil {
			l.writer = os.Stderr
		}

		// otherwise init output with conservative defaults
		l.logger = &logrus.Logger{
			Out:       l.writer,
			Formatter: new(TextFormatter),
			Hooks:     make(LevelHooks),
			Level:     DebugLevel,
			ExitFunc:  closer.Exit,
		}

		l.entry = l.logger.WithContext(context.Background())
		l.reloadStackTraceCache()
		l.addDefaultHooks()
		l.mux = new(sync.Mutex)
		l.initDone = true
	})
}

const defaultStackSearchOffset = 1

// reloadStackTraceCache allows to reload the stack trace reporter with new offset,
// allowing to wrap suplogger into other funcs.
func (l *suplogger) reloadStackTraceCache() {
	l.stack = stackcache.New(defaultStackSearchOffset, l.stackTraceOffset, "github.com/xlab/suplog")
}

// addDefaultHooks initializes default hooks and additional hooks
// based on the environment setup.
func (l *suplogger) addDefaultHooks() {
	// new logger with same out and formatter, but no hooks.
	// used to avoid hooking a hooka-roo from hooks,
	// that hits a mutex in the same logrus entry.
	hookLogger := NewLogger(l.logger.Out, l.logger.Formatter)

	l.logger.AddHook(debugHook.NewHook(hookLogger, nil))

	if isTrue(os.Getenv("LOG_BLOB_ENABLED")) {
		l.logger.AddHook(blobHook.NewHook(hookLogger, nil))
	}

	if isTrue(os.Getenv("LOG_BUGSNAG_ENABLED")) {
		l.logger.AddHook(bugsnagHook.NewHook(hookLogger, nil))
	}
}

// Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Error, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (l *suplogger) WithField(key string, value interface{}) Logger {
	l.initOnce()

	outCopy := l.copy()
	outCopy.entry = l.entry.WithField(key, value)

	return outCopy
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (l *suplogger) WithFields(fields Fields) Logger {
	l.initOnce()
	outCopy := l.copy()
	outCopy.entry = l.entry.WithFields(fields)

	return outCopy
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (l *suplogger) WithError(err error) Logger {
	l.initOnce()
	outCopy := l.copy()
	outCopy.entry = l.entry.WithError(err)

	return outCopy
}

// Add a context to the log entry.
func (l *suplogger) WithContext(ctx context.Context) Logger {
	l.initOnce()
	outCopy := l.copy()
	outCopy.entry = l.entry.WithContext(ctx)

	return outCopy
}

// Overrides the time of the log entry.
func (l *suplogger) WithTime(t time.Time) Logger {
	l.initOnce()
	outCopy := l.copy()
	outCopy.entry = l.entry.WithTime(t)

	return outCopy
}

func (l *suplogger) Logf(level Level, format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(level, format, args...)
}

func (l *suplogger) Tracef(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(TraceLevel, format, args...)
}

func (l *suplogger) Debugf(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(DebugLevel, format, args...)
}

func (l *suplogger) Infof(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(InfoLevel, format, args...)
}

func (l *suplogger) Printf(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Printf(format, args...)
}

func (l *suplogger) Warningf(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(WarnLevel, format, args...)
}

func (l *suplogger) Errorf(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(ErrorLevel, format, args...)
}

func (l *suplogger) Fatalf(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(FatalLevel, format, args...)
	l.logger.Exit(1)
}

func (l *suplogger) Panicf(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(PanicLevel, format, args...)
}

func (l *suplogger) Log(level Level, args ...interface{}) {
	l.initOnce()
	l.entry.Log(level, args...)
}

func (l *suplogger) Trace(args ...interface{}) {
	l.initOnce()
	l.entry.Log(TraceLevel, args...)
}

func (l *suplogger) Info(args ...interface{}) {
	l.initOnce()
	l.entry.Log(InfoLevel, args...)
}

func (l *suplogger) Print(args ...interface{}) {
	l.initOnce()
	l.entry.Print(args...)
}

func (l *suplogger) Fatal(args ...interface{}) {
	l.initOnce()
	l.entry.Log(FatalLevel, args...)
	l.logger.Exit(1)
}

func (l *suplogger) Panic(args ...interface{}) {
	l.initOnce()
	l.entry.Log(PanicLevel, args...)
}

func (l *suplogger) Logln(level Level, args ...interface{}) {
	l.initOnce()
	l.entry.Logln(level, args...)
}

func (l *suplogger) Traceln(args ...interface{}) {
	l.initOnce()
	l.entry.Logln(TraceLevel, args...)
}

func (l *suplogger) Debugln(args ...interface{}) {
	l.initOnce()
	l.entry.Logln(DebugLevel, args...)
}

func (l *suplogger) Infoln(args ...interface{}) {
	l.initOnce()
	l.entry.Logln(InfoLevel, args...)
}

func (l *suplogger) Println(args ...interface{}) {
	l.initOnce()
	l.entry.Println(args...)
}

func (l *suplogger) Warningln(args ...interface{}) {
	l.initOnce()
	l.entry.Logln(WarnLevel, args...)
}

func (l *suplogger) Errorln(args ...interface{}) {
	l.initOnce()
	l.entry.Logln(ErrorLevel, args...)
}

func (l *suplogger) Fatalln(args ...interface{}) {
	l.initOnce()
	l.entry.Logln(FatalLevel, args...)
	l.logger.Exit(1)
}

func (l *suplogger) Debug(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(DebugLevel, format, args...)
}

func (l *suplogger) Notification(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(InfoLevel, format, args...)
}

func (l *suplogger) Success(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(InfoLevel, format, args...)
}

func (l *suplogger) Warning(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(WarnLevel, format, args...)
}

func (l *suplogger) Error(format string, args ...interface{}) {
	l.initOnce()
	l.entry.Logf(ErrorLevel, format, args...)
}

func (l *suplogger) Panicln(args ...interface{}) {
	l.initOnce()
	l.entry.Logln(PanicLevel, args...)
}

// SetLevel sets the logger level.
func (l *suplogger) SetLevel(level Level) {
	l.initOnce()
	l.logger.SetLevel(level)
}

// GetLevel returns the logger level.
func (l *suplogger) GetLevel() Level {
	l.initOnce()
	return l.logger.GetLevel()
}

// AddHook adds a hook to the logger hooks.
func (l *suplogger) AddHook(hook Hook) {
	l.initOnce()
	l.logger.AddHook(hook)
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (l *suplogger) IsLevelEnabled(level Level) bool {
	l.initOnce()
	return l.logger.IsLevelEnabled(level)
}

// SetFormatter sets the logger formatter.
func (l *suplogger) SetFormatter(formatter Formatter) {
	l.initOnce()
	l.logger.SetFormatter(formatter)
}

// SetFormatter sets the logger formatter.
func (l *suplogger) SetStackTraceOffset(offset int) {
	l.initOnce()
	l.stackTraceOffset = offset
	l.reloadStackTraceCache()
}

// SetOutput sets the logger suplog.
func (l *suplogger) SetOutput(output io.Writer) {
	l.initOnce()
	l.logger.SetOutput(output)
}

// ReplaceHooks replaces the logger hooks and returns the old ones
func (l *suplogger) ReplaceHooks(hooks LevelHooks) LevelHooks {
	l.initOnce()
	return l.logger.ReplaceHooks(hooks)
}

// Close effectively closes output, closing the underlying writer
// if it implements io.WriteCloser.
func (l *suplogger) Close() (err error) {
	// bail out if already closed
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.closed {
		return
	}

	l.closed = true

	// try to close only WriteClosers
	if outCloser, ok := l.writer.(io.WriteCloser); ok {
		return outCloser.Close()
	}

	return
}

// CallerName returns caller function name.
func (l *suplogger) CallerName() string {
	l.initOnce()
	caller := l.stack.GetCaller()
	parts := strings.Split(caller.Function, "/")
	nameParts := strings.Split(parts[len(parts)-1], ".")

	return nameParts[len(nameParts)-1]
}

func isTrue(v string) bool {
	switch strings.ToLower(v) {
	case "1", "true", "y":
		return true
	}

	return false
}

// copy allows to construct an suplogger copy with new entry.
func (l *suplogger) copy() *suplogger {
	return &suplogger{
		writer:   l.writer,
		logger:   l.logger,
		stack:    l.stack,
		mux:      l.mux,
		initDone: l.initDone,
		closed:   l.closed,
	}
}
