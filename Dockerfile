# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application for Firestore-only deployment
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o lookie ./cmd/lookie

# Runtime stage
FROM alpine:latest

# Install ca-certificates only (no SQLite needed for Firestore)
RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S lookie && \
    adduser -S lookie -u 1001 -G lookie

# Copy binary and configuration
COPY --from=builder /app/lookie .
COPY --from=builder /app/config.production.yaml ./config.yaml

# Change ownership
RUN chown -R lookie:lookie /app

# Switch to non-root user
USER lookie

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default command
CMD ["./lookie"]