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
COPY --from=builder /app/books /bin/books
CMD ["/bin/books"]
