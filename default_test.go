package suplog_test

import (
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	. "github.com/xlab/suplog"
	bugsnagHook "github.com/xlab/suplog/hooks/bugsnag"
	debugHook "github.com/xlab/suplog/hooks/debug"
	"github.com/xlab/suplog/wrapped-test"
)

func TestAll(t *testing.T) {
	Print("This is an example basic message")
	Success("This is an example success message")
	Warning("This is an example warning message")
	Error("This is an example error message")
	Debug("This is an example debug message")
	NewLogger(os.Stderr, nil, debugHook.NewHook(DefaultLogger, nil)).
		Debug("Debug message from non-default suplogger")

	// Test logger wrapping with StackTraceOffset

	logWithOffset := NewLogger(os.Stderr, nil,
		debugHook.NewHook(DefaultLogger, &debugHook.HookOptions{
			StackTraceOffset: 1,
		}),
		bugsnagHook.NewHook(DefaultLogger, &bugsnagHook.HookOptions{
			StackTraceOffset: 1,
		}))
	logWithOffset.(LoggerConfigurator).SetStackTraceOffset(1)

	wrapped.NewTestWrapper(logWithOffset).ErrorText("This is an example error message from wrapped logger")
	errWrapped := errors.New("This is an example wrapped error message from wrapped logger")
	wrapped.NewTestWrapper(logWithOffset).ErrorWrapped(errWrapped)
	wrapped.NewTestWrapper(logWithOffset).DebugText("This is an example debug message from wrapped logger")

	time.Sleep(time.Second)
}

func ExamplePrint() {
	Print("Hello world!")
}

func ExamplePrintf() {
	Printf("Hello world! My name is %s.", "Loggy")
}

func ExampleNotification() {
	Notification("Hello world! My name is %s.", "Loggy")
}

func ExampleSuccess() {
	Success("Hello world! My name is %s.", "Loggy")
}

func ExampleWarning() {
	Warning("Hello world! My name is %s.", "Loggy")
}

func ExampleError() {
	Error("Hello world! My name is %s.", "Loggy")
}

func ExampleDebug() {
	Debug("Hello world! My name is %s.", "Loggy")
}
