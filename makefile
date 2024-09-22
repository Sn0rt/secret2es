BINARY_CLI := secret2es
BINARY_SERVER := secret2es-server
GOARCH := amd64

PLATFORMS := linux darwin windows
os = $(word 1, $@)

.PHONY: all build clean test help $(PLATFORMS)

all: build

build: build-cli build-server

build-cli:
	@echo "Building CLI tool..."
	@go build -o bin/$(BINARY_CLI) ./cmd/cli

build-server:
	@echo "Building HTTP server..."
	@go build -o bin/$(BINARY_SERVER) ./cmd/server

clean:
	@echo "Cleaning up..."
	@rm -rf bin/*

test:
	@go test -v ./...

$(PLATFORMS):
	@echo "Building for $(os)..."
	@GOOS=$(os) GOARCH=$(GOARCH) go build -o bin/$(BINARY_CLI)-$(os)-$(GOARCH) ./cmd/cli
	@GOOS=$(os) GOARCH=$(GOARCH) go build -o bin/$(BINARY_SERVER)-$(os)-$(GOARCH) ./cmd/server

help:
	@echo "Available commands:"
	@echo "  make build       - Build both CLI and server binaries for the current platform"
	@echo "  make build-cli   - Build only the CLI binary"
	@echo "  make build-server- Build only the server binary"
	@echo "  make clean       - Remove all binaries from the bin directory"
	@echo "  make test        - Run all tests"
	@echo "  make linux       - Build binaries for Linux"
	@echo "  make darwin      - Build binaries for macOS"
	@echo "  make windows     - Build binaries for Windows"