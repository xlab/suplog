package bugsnag

import (
	"fmt"
	"os"
	"strings"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/sirupsen/logrus"

	"github.com/xlab/suplog/stackcache"
)

// HookOptions allows to set additional Hook options.
type HookOptions struct {
	// Levels enables this hook for all listed levels.
	Levels []logrus.Level
	// StackTraceOffset allows to wrap logger into greater stack depth and still
	// get reports on accurate positions.
	StackTraceOffset int

	Env               string
	AppVersion        string
	BugsnagAPIKey     string
	BugsnagEnabledEnv []string
	BugsnagPackages   []string
	BugsnagCatchAll   bool
}

func checkHookOptions(opt *HookOptions) *HookOptions {
	if opt == nil {
		opt = &HookOptions{}
	}

	if len(opt.Levels) == 0 {
		opt.Levels = []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		}
	}

	if len(opt.Env) == 0 {
		opt.Env = os.Getenv("APP_ENV")
		if len(opt.Env) == 0 {
			opt.Env = "local"
		}
	}

	if len(opt.AppVersion) == 0 {
		opt.AppVersion = os.Getenv("APP_VERSION")
	}

	if len(opt.BugsnagAPIKey) == 0 {
		opt.BugsnagAPIKey = os.Getenv("LOG_BUGSNAG_KEY")
	}

	if !opt.BugsnagCatchAll {
		opt.BugsnagCatchAll = toBool(os.Getenv("LOG_BUGSNAG_CATCH_ALL"))
	}

	if len(opt.BugsnagEnabledEnv) == 0 {
		opt.BugsnagEnabledEnv = []string{
			"prod",
			"staging",
			"test",
		}
	}

	if len(opt.BugsnagPackages) == 0 {
		opt.BugsnagPackages = []string{
			"main",
			"github.com/xlab/suplog/*",
		}
	}

	return opt
}

type RootLogger interface {
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Printf(format string, args ...interface{})
}

const defaultStackSearchOffset = 6

// NewHook initializes a new logrus.Hook using provided params and options.
// Provide a root logger to print any errors occuring during the plugin init.
func NewHook(logger RootLogger, opt *HookOptions) logrus.Hook {
	opt = checkHookOptions(opt)

	var panicHandler func() = nil
	if !opt.BugsnagCatchAll {
		panicHandler = func() {}
	}

	return &hook{
		opt:    opt,
		logger: logger,
		stack:  stackcache.New(defaultStackSearchOffset, opt.StackTraceOffset, "github.com/xlab/suplog"),
		notifier: bugsnag.New(bugsnag.Configuration{
			APIKey:              opt.BugsnagAPIKey,
			ReleaseStage:        opt.Env,
			ProjectPackages:     opt.BugsnagPackages,
			AppVersion:          opt.AppVersion,
			NotifyReleaseStages: opt.BugsnagEnabledEnv,
			PanicHandler:        panicHandler,
			Logger:              logger,
		}),
	}
}

type hook struct {
	opt      *HookOptions
	logger   RootLogger
	stack    stackcache.StackCache
	notifier *bugsnag.Notifier
}

func (h *hook) Levels() []logrus.Level {
	return h.opt.Levels
}

func (h *hook) Fire(e *logrus.Entry) error {
	var (
		err        ErrorWithStackFrames
		errContext bugsnag.Context
	)

	// check if we have error in fields
	if withErr, ok := e.Data["error"].(error); ok {
		// check if that error has stack (was wrapped at some point)
		if withStack, ok := withErr.(ErrorWithStackFrames); ok {
			// use this error to report, with its original stack
			err = withStack
			errContext.String = e.Message
		} else if stackTracer, ok := withErr.(pkgErrorsStackTracer); ok {
			// the error is pkg/errors wrapped error, try to parse it
			var parsingErr error

			stackTrace := stackTracer.StackTrace()
			err, parsingErr = newErrorWithPkgErrorsStackTrace(withErr, stackTrace)
			if parsingErr != nil {
				// no stack with error (parsing failure), wrap it
				stackFrames := h.stack.GetStackFrames()
				err = newErrorWithStackFrames(withErr, stackFrames)
			}

			errContext.String = e.Message
		} else {
			// no stack with error, wrap it
			stackFrames := h.stack.GetStackFrames()
			err = newErrorWithStackFrames(withErr, stackFrames)
			errContext.String = e.Message
		}
	} else {
		// no error within fields, construct new one from log message
		stackFrames := h.stack.GetStackFrames()
		err = newErrorWithStackFrames(fmt.Errorf("%s", e.Message), stackFrames)
	}

	var (
		needSync = false
		severity = bugsnag.SeverityInfo
	)

	switch e.Level {
	case logrus.WarnLevel:
		severity = bugsnag.SeverityWarning
	case logrus.ErrorLevel:
		severity = bugsnag.SeverityError
	case logrus.FatalLevel, logrus.PanicLevel:
		severity = bugsnag.SeverityError
		needSync = true
	}

	userData := captureUserMeta(e.Data)
	metaData := fieldsToMetaData(e.Data)

	if len(errContext.String) > 0 {
		_ = h.notifier.NotifySync(err, needSync, severity, metaData, userData, errContext)
		return nil
	}

	_ = h.notifier.NotifySync(err, needSync, severity, metaData, userData)

	return nil
}

func captureUserMeta(fields logrus.Fields) (user bugsnag.User) {
	if userID, ok := fields["@user.id"].(string); ok {
		user.Id = userID

		delete(fields, "@user.id")
	}

	if userName, ok := fields["@user.name"].(string); ok {
		user.Name = userName

		delete(fields, "@user.name")
	}

	if userEmail, ok := fields["@user.email"].(string); ok {
		user.Email = userEmail

		delete(fields, "@user.email")
	}

	return user
}

func fieldsToMetaData(fields logrus.Fields) bugsnag.MetaData {
	if len(fields) == 0 {
		return bugsnag.MetaData{}
	}

	fieldsMap := make(map[string]interface{}, len(fields))

	for field, value := range fields {
		switch field {
		case "blob", "error":
			continue
		}

		fieldsMap[field] = value
	}

	return bugsnag.MetaData{
		"Fields": fieldsMap,
	}
}

func toBool(s string) bool {
	switch strings.ToLower(s) {
	case "true", "1", "t", "yes":
		return true
	default:
		return false
	}
}
