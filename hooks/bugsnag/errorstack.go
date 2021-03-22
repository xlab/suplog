package bugsnag

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/bugsnag/bugsnag-go/errors"
	pkgerrors "github.com/pkg/errors"

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

type pkgErrorsStackTracer interface {
	StackTrace() pkgerrors.StackTrace
}

func newErrorWithPkgErrorsStackTrace(err error, stackTrace pkgerrors.StackTrace) (ErrorWithStackFrames, error) {
	if err == nil {
		return nil, nil
	}

	e := &errorWithStack{
		orig:   err,
		frames: make([]errors.StackFrame, len(stackTrace)),
	}

	for i, frame := range stackTrace {
		var (
			fnName         string
			fileName       string
			fileLineNumber int
		)

		frameText, parseErr := frame.MarshalText()
		if parseErr != nil {
			return nil, parseErr
		}

		parts := strings.Split(string(frameText), " ")
		if len(parts) != 2 {
			parseErr = fmt.Errorf("frame text partial read: not enough parts")
			return nil, parseErr
		}
		fnName = parts[0]

		lineNumIdx := strings.LastIndexByte(parts[1], ':')
		if lineNumIdx < 0 || lineNumIdx+1 >= len(parts[1]) {
			parseErr = fmt.Errorf("frame text partial read: no file line delim in %s", parts[1])
			return nil, parseErr
		}

		fileLineNumber, parseErr = strconv.Atoi(parts[1][lineNumIdx+1:])
		if parseErr != nil {
			parseErr = fmt.Errorf("failed to parse line num %s", parseErr)
			return nil, parseErr
		}
		fileName = parts[1][:lineNumIdx]

		e.frames[i] = errors.StackFrame{
			File:           limitPath(fileName, 3),
			LineNumber:     fileLineNumber,
			Name:           fnName,
			Package:        stackcache.GetPackageName(fnName),
			ProgramCounter: uintptr(frame),
		}
	}

	return e, nil
}
