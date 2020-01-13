package stackcache

import (
	"runtime"
	"strings"
	"sync"
)

type StackCache interface {
	GetCaller() runtime.Frame
	GetStackFrames() []runtime.Frame
}

// New creates a new stack cache for effectively traversing runtime frame stack.
// The traverser will start at pcOffset and move until not exited from internal
// packages of the output library.
func New(pcOffset int, breakpointPackage string) StackCache {
	return &stackCache{
		minimumCallerDepth: pcOffset,
		maximumCallerDepth: 25,
		breakpointPackage:  breakpointPackage,
	}
}

type stackCache struct {
	// breakpointPackage is a package name that stack traverser will seek,
	// so it could ignore frames upon finding the first frame after that package.
	breakpointPackage string

	offsetOnce         sync.Once
	offsetIsSet        bool
	minimumCallerDepth int
	maximumCallerDepth int
}

// pkgNameTesting is the package of testing.tRunner, in case if
// the top-most package that calls output library is a default test runner.
const pkgNameTesting = "testing"

// GetCaller retrieves the name of the first function from a non-internal package.
// That would be our caller.
func (c *stackCache) GetCaller() runtime.Frame {
	pcs := make([]uintptr, c.maximumCallerDepth)
	depth := runtime.Callers(c.minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	var (
		offset      int
		latestFrame runtime.Frame
	)

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := GetPackageName(f.Function)

		if !c.offsetIsSet {
			if pkg == c.breakpointPackage {
				c.offsetOnce.Do(func() {
					c.minimumCallerDepth += offset
					c.offsetIsSet = true
				})
			}
		} else if pkg != c.breakpointPackage {
			if pkg == pkgNameTesting {
				break
			}
			latestFrame = f
			break
		}

		latestFrame = f
		offset++
	}

	return latestFrame
}

// GetStackFrames retrieves the full stack since first non-internal package.
func (c *stackCache) GetStackFrames() []runtime.Frame {
	pcs := make([]uintptr, c.maximumCallerDepth)
	depth := runtime.Callers(c.minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	usefulStackFrames := make([]runtime.Frame, 0, depth)

	var (
		offset      int
		latestFrame runtime.Frame
		latestPkg   string
	)

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := GetPackageName(f.Function)

		if !c.offsetIsSet {
			if pkg == c.breakpointPackage {
				c.offsetOnce.Do(func() {
					c.minimumCallerDepth += offset
					c.offsetIsSet = true
				})
			}
		} else if pkg != c.breakpointPackage {
			if pkg == pkgNameTesting && latestPkg == c.breakpointPackage {
				usefulStackFrames = append(usefulStackFrames, latestFrame)
			}
			usefulStackFrames = append(usefulStackFrames, f)
			latestPkg = pkg
			continue
		}

		latestFrame = f
		latestPkg = pkg
		offset++
	}

	return usefulStackFrames
}

// GetPackageName reduces a fully qualified function name to the package name
// This function is from logrus internals.
func GetPackageName(path string) string {
	for {
		lastPeriod := strings.LastIndex(path, ".")
		lastSlash := strings.LastIndex(path, "/")

		if lastPeriod > lastSlash {
			path = path[:lastPeriod]
		} else {
			break
		}
	}

	return path
}
