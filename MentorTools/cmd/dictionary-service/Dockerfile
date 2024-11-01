# cmd/dictionary-service-service/Dockerfile

# Stage 1: Build the application
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/dictionary-service ./cmd/dictionary-service
COPY internal/dictionary-service ./internal/dictionary-service
COPY pkg ./pkg

# Build the binary
RUN go build -o dictionary-service ./cmd/dictionary-service/main.go

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/dictionary-service .

# Указываем внутренний порт приложения
EXPOSE 8080

# Команда для запуска приложения
CMD ["./dictionary-service"]