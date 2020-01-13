package blob

import (
	"bytes"
	"io"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Remote provides Amazon S3 compatible bucket access methods.
type S3Remote interface {
	CheckAccess(key string) error
	PutObject(key string, r io.Reader, meta map[string]string) (*S3Spec, error)
}

func NewS3Remote(accoutID, secretKey, endpoint, region, bucket string) (s3Client S3Remote, err error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accoutID, secretKey, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	s3Client = &s3Remote{
		bucket: bucket,
		cli:    s3.New(sess),
	}

	return s3Client, nil
}

type S3Spec struct {
	Path      string
	Key       string
	Body      io.ReadCloser
	ETag      string
	Version   string
	UpdatedAt time.Time
	Meta      map[string]string
	Size      int64
}

type s3Remote struct {
	bucket string
	cli    *s3.S3
}

func (s *s3Remote) CheckAccess(prefix string) error {
	body := []byte(time.Now().UTC().String())
	_, err := s.cli.PutObject(&s3.PutObjectInput{
		Body:        aws.ReadSeekCloser(bytes.NewReader(body)),
		Bucket:      aws.String(s.bucket),
		ContentType: aws.String("text/plain"),
		Key:         aws.String(path.Join(prefix, "_touch")),
	})

	return err
}

func (s *s3Remote) PutObject(key string, r io.Reader, meta map[string]string) (*S3Spec, error) {
	obj, err := s.cli.PutObject(&s3.PutObjectInput{
		Body:     aws.ReadSeekCloser(r),
		Bucket:   aws.String(s.bucket),
		Key:      aws.String(key),
		Metadata: aws.StringMap(meta),
	})
	if err != nil {
		return nil, err
	}

	spec := &S3Spec{
		Key:     key,
		ETag:    aws.StringValue(obj.ETag),
		Version: aws.StringValue(obj.VersionId),
		Meta:    meta,
	}

	return spec, err
}
