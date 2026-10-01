# Multi-stage Dockerfile for Production
# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install CA SSL certificates for secure HTTPS connections (LinkedIn & Gemini APIs)
RUN apk add --no-cache ca-certificates

# Copy Go module dependency definitions
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build optimized static binary without CGO
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/git-to-feed ./cmd/api

# Final Stage: Minimalist runtime image
FROM alpine:3.19

WORKDIR /app

# Copy ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy executable binary
COPY --from=builder /app/git-to-feed /app/git-to-feed

# Copy SQL migrations folder
COPY --from=builder /app/migrations /app/migrations

# Expose default application port
EXPOSE 8080

# Unprivileged non-root user for security
RUN adduser -D appuser
USER appuser

# Execution entrypoint
ENTRYPOINT ["/app/git-to-feed"]