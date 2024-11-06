# Этап сборки
FROM golang:1.22-alpine AS builder

# Установка необходимых пакетов
RUN apk update && apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN go build -o subscription-service ./cmd/main.go

# Этап выполнения
FROM alpine:latest

# Установка ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates

# Создание рабочей директории
WORKDIR /root/

# Копирование бинарника из этапа сборки
COPY --from=builder /app/subscription-service .

# Открытие порта приложения
EXPOSE 8080

# Команда запуска приложения
CMD ["./subscription-service"]
