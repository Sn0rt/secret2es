# Build stage
FROM golang:1.20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
ARG VERSION
ARG BUILD_TIME
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X cmd.version=${VERSION} -X cmd.buildTime=${BUILD_TIME}" -o sercert2extsecret

# Final stage
FROM alpine:3.14

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/sercert2extsecret .

# Command to run the executable
ENTRYPOINT ["./sercert2es"]