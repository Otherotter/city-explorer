FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o collector ./services/collector

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/collector .
EXPOSE 8081
CMD ["./collector"]