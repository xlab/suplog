module github.com/xlab/suplog

go 1.14

require (
	github.com/aws/aws-sdk-go v1.25.16
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bugsnag/bugsnag-go v1.5.3
	github.com/kr/pretty v0.1.0 // indirect
	github.com/oklog/ulid v1.3.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/xlab/closer v0.0.0-20190328110542-03326addb7c2
	golang.org/x/net v0.0.0-20190923162816-aa69164e4478 // indirect
	golang.org/x/sys v0.0.0-20200803210538-64077c9b5642 // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace github.com/bugsnag/bugsnag-go => ./hooks/bugsnag/bugsnag-go
