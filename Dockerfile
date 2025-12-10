FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git dan ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go files first
COPY go.mod go.sum ./

# Download dependencies dengan verbose untuk debugging
RUN go mod download -x

# Copy semua source code
COPY . .

# Tidy module dan build dengan verbose untuk debugging
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -v -o api ./cmd/api

# Production image
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy binary dan files yang dibutuhkan
COPY --from=builder /app/api .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/configs ./configs

EXPOSE 9001

CMD ["./api"]