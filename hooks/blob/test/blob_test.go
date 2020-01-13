package blob

import (
	"os"
	"testing"
	"time"

	blobHook "github.com/xlab/suplog/hooks/blob"

	"github.com/xlab/suplog"
)

func TestBlobHook(t *testing.T) {
	testBlob := []byte(`Son agreed others exeter period myself few yet nature.
		Mention mr manners opinion if garrets enabled.
		To an occasional dissimilar impossible sentiments.
		Do fortune account written prepare invited no passage.
		Garrets use ten you the weather ferrars venture friends.
		Solid visit seems again you nor all.`)

	opts := &blobHook.HookOptions{
		Env: "test",
	}

	out := suplog.NewLogger(
		os.Stderr,
		new(suplog.TextFormatter),
		blobHook.NewHook(suplog.DefaultLogger, opts),
	)
	ts := time.Now()

	out.WithField("blob", testBlob).Infoln("test is running, trying to submit blob")
	out.Debug("test done in %s", time.Since(ts))
}
