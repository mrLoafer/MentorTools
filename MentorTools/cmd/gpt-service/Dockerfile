# cmd/gpt-service-service/Dockerfile

# Stage 1: Build the application
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/gpt-service ./cmd/gpt-service
COPY internal/gpt-service ./internal/gpt-service
COPY pkg ./pkg

# Build the binary
RUN go build -o gpt-service ./cmd/gpt-service/main.go

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/gpt-service .

# Указываем внутренний порт приложения
EXPOSE 8080

# Команда для запуска приложения
CMD ["./gpt-service"]