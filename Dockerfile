# Stage 1: Build
FROM golang:1.23.0-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["/app/main"]