FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 8080 9090

CMD ["./main"]