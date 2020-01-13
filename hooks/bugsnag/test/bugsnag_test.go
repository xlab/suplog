package bugsnag

import (
	"errors"
	"os"
	"testing"
	"time"

	bugsnagHook "github.com/xlab/suplog/hooks/bugsnag"

	"github.com/xlab/suplog"
)

func TestBugsnagHook(t *testing.T) {
	opts := &bugsnagHook.HookOptions{
		Env:        "test",
		AppVersion: "magic_horse",
	}

	out := suplog.NewLogger(
		os.Stderr,
		new(suplog.TextFormatter),
		bugsnagHook.NewHook(suplog.DefaultLogger, opts),
	)

	out.Info("test has started")
	out.Error("1) some fake error as text")
	suplog.Warning("2) also default suplogger with enabled env")
	out.WithError(errors.New("some fake error")).WithFields(suplog.Fields{
		"@user.name": "Max",
	}).Error("3) with fields and error, also meta")

	time.Sleep(time.Second)
	out.Debug("test done")
}
