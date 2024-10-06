# Сборка приложения
FROM golang:1.22 AS builder

WORKDIR /msg-bakend

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем Go-приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Финальный этап
FROM alpine:latest

WORKDIR /root

# Копируем исполняемый файл из сборочного контейнера
COPY --from=builder /msg-bakend/main .

# Открываем порт для приложения (если нужно)
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
