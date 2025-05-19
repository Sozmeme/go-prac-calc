# Используем официальный образ Go
FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go build -o main .

# Открываем порты которые использует приложение
EXPOSE 8080 9090

CMD ["./main"]