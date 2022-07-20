module github.com/xlab/suplog/hooks/bugsnag

go 1.16

require (
	github.com/bugsnag/bugsnag-go v1.5.3
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	github.com/xlab/suplog v1.4.1
)

replace github.com/bugsnag/bugsnag-go => ./bugsnag-go
