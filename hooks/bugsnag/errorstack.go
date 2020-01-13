package bugsnag

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bugsnag/bugsnag-go/errors"

	"github.com/xlab/suplog/stackcache"
)

// ErrorWithStackFrames is acceptable by Bugsnag, it provides a convenient
// way to construct a custom stack (captured with stackcache).
type ErrorWithStackFrames interface {
	Error() string
	StackFrames() []errors.StackFrame
}

type errorWithStack struct {
	orig   error
	frames []errors.StackFrame
}

var _ errors.ErrorWithStackFrames = &errorWithStack{}

func newErrorWithStackFrames(err error, stackFrames []runtime.Frame) ErrorWithStackFrames {
	if err == nil {
		return nil
	}

	e := &errorWithStack{
		orig:   err,
		frames: make([]errors.StackFrame, len(stackFrames)),
	}

	for i, frame := range stackFrames {
		e.frames[i] = errors.StackFrame{
			File:           limitPath(frame.File, 3),
			LineNumber:     frame.Line,
			Name:           frame.Function,
			Package:        stackcache.GetPackageName(frame.Function),
			ProgramCounter: frame.PC,
		}
	}

	return e
}

func (e *errorWithStack) Error() string {
	return e.orig.Error()
}

func (e *errorWithStack) StackFrames() []errors.StackFrame {
	return e.frames
}

func limitPath(path string, n int) string {
	if n <= 0 {
		return path
	}

	pathParts := strings.Split(path, string(filepath.Separator))
	if len(pathParts) > n {
		pathParts = pathParts[len(pathParts)-n:]
	}

	return filepath.Join(pathParts...)
}
