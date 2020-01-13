module github.com/xlab/suplog

go 1.13

require (
	github.com/aws/aws-sdk-go v1.25.16
	github.com/bugsnag/bugsnag-go v1.5.3
	github.com/golangci/golangci-lint v1.22.2 // indirect
	github.com/oklog/ulid v1.3.1
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/xlab/closer v0.0.0-20190328110542-03326addb7c2
)

replace github.com/bugsnag/bugsnag-go => ./hooks/bugsnag/bugsnag-go
