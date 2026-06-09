# Stage 1: Build the Go binary
FROM golang:alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Compile optimized static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api-server cmd/api-server/main.go

# Stage 2: Minimal runner stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /app/api-server .

# Create configuration folder mapping to prevent missing file warnings
RUN mkdir -p internal/config && touch internal/config/.env

EXPOSE 1323

ENTRYPOINT ["./api-server"]
