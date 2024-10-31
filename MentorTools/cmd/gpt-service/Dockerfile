# cmd/gpt-service/Dockerfile

# Stage 1: Build the application
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Copy the config.yaml file
COPY MentorTools/pkg/config/config.yaml /app/config/config.yaml

# Собираем бинарный файл
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