// Package wrapped provides testing primitives for validating stack discovery by suplog.
package wrapped

import (
	log "github.com/xlab/suplog"
)

// TestWrapper is a primitive to test error logging from an extra
// package used as an adapter for logging.
type TestWrapper interface {
	ErrorText(str string)
	ErrorWrapped(err error)
	DebugText(str string)
}

// NewTestWrapper returns a new test wrapper with initialized logger.
func NewTestWrapper(logger log.Logger) TestWrapper {
	return &testWrapper{
		logger: logger,
	}
}

type testWrapper struct {
	logger log.Logger
}

// ErrorText logs just error text, meaning that stacktrace will be captured from there
func (e *testWrapper) ErrorText(str string) {
	e.logger.Error(str)
}

// ErrorWrapped accepts an error wrapped with stacktrace, like github.com/pkg/errors
func (e *testWrapper) ErrorWrapped(err error) {
	e.logger.WithError(err).Error("error wrapped")
}

// DebugText captures strack trace in the log report, if debug hook is enabled.
func (e *testWrapper) DebugText(str string) {
	e.logger.Debugln(str)
}
