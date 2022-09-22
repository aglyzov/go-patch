TEST_TARGETS := $(shell go list ./pkg/...)
# Go source files, ignore vendor directory
SRC := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

integration-test:
	export TEST_INTEGRATION=true && \
	go test $(TEST_TARGETS) -v -failfast -count=1 -race -parallel=6

unit-test:
	go test $(TEST_TARGETS) -v -failfast -count=1 -race -parallel=6

format:
	@gofmt -l -w $(SRC)
	@goimports -w -e -local github.com/itsrever/go-patch

install-lint-ubuntu:
	echo Installing yamlint golangci-lint...
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.44.2
	golangci-lint --version

install-lint-macos:
	brew install golangci-lint

lint: format
	@golangci-lint -v --timeout=600s --skip-dirs=docs run