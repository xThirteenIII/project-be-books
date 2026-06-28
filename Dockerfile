# Go Builder, avoids external building
# Use 1.24 stable version
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o books ./cmd/api

# Minimal image
FROM debian:stable-slim

# Handle ca-certificates or the Go app can't talk with Gutendex
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app

COPY --from=builder /app/books /bin/books
CMD ["/bin/books"]
