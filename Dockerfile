FROM golang:1.20-alpine AS builder

ENV APP_NAME=run_app

RUN apk add --no-cache gcc musl-dev sqlite-dev make

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 go build -o ./$APP_NAME cmd/main.go

RUN chmod +x docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]