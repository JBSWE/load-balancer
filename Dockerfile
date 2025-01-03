FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

WORKDIR /app/cmd/server

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/cmd/server/main .

EXPOSE 8080

CMD ["./main"]
