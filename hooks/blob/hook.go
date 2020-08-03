package blob

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// HookOptions allows to set additional Hook options.
type HookOptions struct {
	Env               string
	BlobStoreURL      string
	BlobStoreAccount  string
	BlobStoreKey      string
	BlobStoreEndpoint string
	BlobStoreRegion   string
	BlobStoreBucket   string
	BlobRetentionTTL  time.Duration
	BlobEnabledEnv    map[string]bool
}

// DefaultRetentionTTL is currently set to be 1 month.
const DefaultRenentionTTL = 30 * 24 * time.Hour

func checkHookOptions(opt *HookOptions) *HookOptions {
	if opt == nil {
		opt = &HookOptions{}
	}

	if len(opt.Env) == 0 {
		opt.Env = os.Getenv("APP_ENV")
		if len(opt.Env) == 0 {
			opt.Env = "local"
		}
	}

	if len(opt.BlobStoreURL) == 0 {
		opt.BlobStoreURL = os.Getenv("LOG_BLOB_STORE_URL")
	}

	if len(opt.BlobStoreAccount) == 0 {
		opt.BlobStoreAccount = os.Getenv("LOG_BLOB_STORE_ACCOUNT")
	}

	if len(opt.BlobStoreKey) == 0 {
		opt.BlobStoreKey = os.Getenv("LOG_BLOB_STORE_KEY")
	}

	if len(opt.BlobStoreEndpoint) == 0 {
		opt.BlobStoreEndpoint = os.Getenv("LOG_BLOB_STORE_ENDPOINT")
	}

	if len(opt.BlobStoreRegion) == 0 {
		opt.BlobStoreRegion = os.Getenv("LOG_BLOB_STORE_REGION")
	}

	if len(opt.BlobStoreBucket) == 0 {
		opt.BlobStoreBucket = os.Getenv("LOG_BLOB_STORE_BUCKET")
	}

	if opt.BlobRetentionTTL == 0 {
		opt.BlobRetentionTTL = DefaultRenentionTTL
	}

	if len(opt.BlobEnabledEnv) == 0 {
		opt.BlobEnabledEnv = map[string]bool{
			"prod":    true,
			"staging": true,
			"test":    true,
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

// NewHook initializes a new suplog.Hook using provided params and options.
// Provide a root logger to print any errors occuring during the plugin init.
func NewHook(logger RootLogger, opt *HookOptions) logrus.Hook {
	h := &hook{
		logger: logger,
		opt:    checkHookOptions(opt),
	}

	if s3Remote, err := NewS3Remote(
		h.opt.BlobStoreAccount,
		h.opt.BlobStoreKey,
		h.opt.BlobStoreEndpoint,
		h.opt.BlobStoreRegion,
		h.opt.BlobStoreBucket,
	); err != nil {
		logger.Errorf("failed to init S3 session: %+v", err)
		return h
	} else if err = s3Remote.CheckAccess(h.opt.Env); err != nil {
		logger.Errorf("failed to verify S3 remote access: %+v", err)
		return h
	} else {
		h.s3Remote = s3Remote
	}

	return h
}

type hook struct {
	opt      *HookOptions
	logger   RootLogger
	s3Remote S3Remote
}

func (h *hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

func (h *hook) Fire(e *logrus.Entry) error {
	blob, hasBlob := e.Data["blob"]
	if !hasBlob {
		return nil
	}

	if h.s3Remote == nil {
		h.logger.Warningf("blob provided but S3 remote is disabled")
		delete(e.Data, "blob")

		return nil
	} else if enabled := h.opt.BlobEnabledEnv[h.opt.Env]; !enabled {
		h.logger.Debugf("blob provided but uploading is disabled in %s", h.opt.Env)
		delete(e.Data, "blob")
		return nil
	}

	var blobPayload []byte
	switch bb := blob.(type) {
	case string:
		blobPayload = []byte(bb)
	case []byte:
		blobPayload = make([]byte, len(bb))
		copy(blobPayload, bb)
	default:
		delete(e.Data, "blob")
		return nil
	}

	blobID := NewBlobID()

	if len(h.opt.BlobStoreURL) > 0 {
		e.Data["blob"] = fmt.Sprintf("%s/%s", h.opt.BlobStoreURL, blobID)
	} else {
		e.Data["blob"] = fmt.Sprintf("%s/%s", h.opt.Env, blobID)
	}

	h.blobUpload(blobID, blobPayload)

	return nil
}

func (h *hook) blobUpload(blobID string, payload []byte) {
	objectKey := filepath.Join(h.opt.Env, blobID)
	_, err := h.s3Remote.PutObject(objectKey, bytes.NewReader(payload), nil)

	if err != nil {
		h.logger.Errorf(
			"failed to upload blob to S3 remote server: key %s in %s: %+v",
			objectKey,
			h.opt.BlobStoreBucket,
			err,
		)
	}
}
