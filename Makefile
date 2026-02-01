.PHONY: build test lint clean install all

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

all: lint test build

build:
	go build $(LDFLAGS) -o jcli .

build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/jcli-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/jcli-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/jcli-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/jcli-darwin-arm64 .

test:
	go test -v -race -coverprofile=coverage.out ./...

test-integration:
	go test -v -tags=integration ./tests/integration/...

lint:
	golangci-lint run ./...

clean:
	rm -f jcli coverage.out
	rm -rf dist/

install: build
	cp jcli $(GOPATH)/bin/

coverage: test
	go tool cover -html=coverage.out -o coverage.html
