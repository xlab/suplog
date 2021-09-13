module github.com/xlab/suplog

go 1.15

require (
	github.com/aws/aws-sdk-go v1.25.16
	github.com/bugsnag/bugsnag-go v1.5.3
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/oklog/ulid v1.3.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/xlab/closer v0.0.0-20190328110542-03326addb7c2
	golang.org/x/net v0.0.0-20201202161906-c7110b5ffcbb // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace github.com/bugsnag/bugsnag-go => ./hooks/bugsnag/bugsnag-go
