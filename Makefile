all:

install: export GO111MODULE=on
install:
	go install github.com/xlab/suplog

lint:
	golangci-lint run --enable-all -D gomnd

test:
	go test
