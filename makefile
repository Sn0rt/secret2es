BINARY_CLI := secret2es
BINARY_SERVER := secret2es-server
DOCKER_REPO=wangguohao/secret2es

# Build information
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

all: build-cli build-server

build: build-cli

build-cli:
	@echo "Building CLI tool..."
	@go build -v -o $(BINARY_CLI) $(LDFLAGS) ./cmd/cli

build-server:
	@echo "Building HTTP server..."
	@go build -v -o $(BINARY_SERVER) $(LDFLAGS) ./cmd/server

print-binary-name:
	@echo $(BINARY_CLI)

test:
	@go test -v ./...

clean:
	@echo "Cleaning up..."
	@rm -rf bin/*

docker-build:
	docker build --build-arg VERSION=$(VERSION) --build-arg BUILD_TIME=$(BUILD_TIME) -t $(DOCKER_REPO):$(VERSION) .

help:
	@echo "Available commands:"
	@echo "  make build       - Build both CLI and server binaries for the current platform"
	@echo "  make build-cli   - Build only the CLI binary"
	@echo "  make build-server- Build only the server binary"
	@echo "  make clean       - Remove all binaries from the bin directory"
	@echo "  make test        - Run tests"

.PHONY: all build build-cli build-server print-binary-name test clean docker-build help