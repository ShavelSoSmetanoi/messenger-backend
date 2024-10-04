# Используем официальный образ Go для сборки
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /msg-bakend/cmd/app/

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект в контейнер
COPY ./cmd/app .

# Собираем Go-приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o msg-bakend .

# Используем минимальный образ для запуска приложения
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированное приложение из образа builder
COPY --from=builder /msg-bakend/msg-bakend .

# Открываем порт для приложения
EXPOSE 8080

# Команда для запуска приложения
CMD ["./msg-bakend"]
