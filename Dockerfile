# --- Stage 1: Build ---
FROM golang:1.22-bullseye AS builder

WORKDIR /go/src/app

# Install dependencies first (better cache)
COPY backend/go.mod backend/go.sum ./
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

# Expose service port
EXPOSE 8081

# Set environment defaults (can be overridden in docker-compose.yml)
ENV DB_HOST=db \
    DB_USER=postgres \
    DB_PASSWORD=example

# Run backend
CMD ["./main"]
