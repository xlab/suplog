package suplog

import (
	"os"
	"testing"
	"time"

	debugHook "github.com/xlab/suplog/hooks/debug"
)

func TestAll(t *testing.T) {
	Print("This is an example basic message")
	Success("This is an example success message")
	Warning("This is an example warning message")
	Error("This is an example error message")
	Debug("This is an example debug message")
	NewLogger(os.Stderr, nil, debugHook.NewHook(DefaultLogger, nil)).
		Debug("Debug message from non-default suplogger")
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
