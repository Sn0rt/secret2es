# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=sercert2es
BINARY_UNIX=$(BINARY_NAME)_unix

# Build information
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

all: test build

build:
	$(GOBUILD) -v -o $(BINARY_NAME) ${LDFLAGS} cmd/secret2es.go

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) ${LDFLAGS} -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) ${LDFLAGS} -o $(BINARY_UNIX) -v

docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

.PHONY: all build test clean run build-linux docker-build