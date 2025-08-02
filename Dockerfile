# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o lookie cmd/lookie/main.go

# Runtime stage
FROM alpine:latest

# Install sqlite and ca-certificates
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S lookie && \
    adduser -S lookie -u 1001 -G lookie

# Create data directory
RUN mkdir -p /app/data && \
    chown -R lookie:lookie /app

# Copy binary and migrations
COPY --from=builder /app/lookie .
COPY --from=builder /app/migrations ./migrations/
COPY --from=builder /app/config.example.yaml ./config.yaml

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