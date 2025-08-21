# --- Stage 1: Build ---
FROM golang:1.24-bullseye AS builder

WORKDIR /go/src/app

# Install dependencies first (better cache)
COPY go.mod go.sum ./
RUN go mod download && go mod tidy

# Copy source
COPY . .

# Build the Go binary
RUN go build -o main ./cmd/main.go


# --- Stage 2: Run ---
FROM debian:bullseye-slim

WORKDIR /app

# Copy built binary
COPY --from=builder /go/src/app/main .

# Copy .env
COPY .env .

# Expose service port
EXPOSE 8081

# Run backend
CMD ["./main"]