# Stage 1: Build the application
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code and config.yaml
COPY . .
COPY pkg/config/config.yaml /app/config/config.yaml

# Build the binary
RUN go build -o auth-service ./cmd/auth-service/main.go

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /app

# Copy binary and config.yaml from builder
COPY --from=builder /app/auth-service .
COPY --from=builder /app/config/config.yaml ./config/config.yaml

# Expose internal app port
EXPOSE 8080

# Run the app
CMD ["./auth-service"]