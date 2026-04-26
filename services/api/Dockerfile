# services/api/Dockerfile

# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api main.go

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]