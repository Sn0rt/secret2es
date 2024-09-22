# Build stage
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /go/src/github.com/sn0rt/secret2es

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build arguments
ARG VERSION
ARG BUILD_TIME

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o secret2es \
    -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}" \
    cmd/cli/main.go

# Final stage
FROM alpine:3.18

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /go/src/github.com/sn0rt/secret2es/secret2es .
