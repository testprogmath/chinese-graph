# Build stage
FROM golang:1.25.8-alpine3.23 AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Debug: Show what was copied
RUN ls -la
RUN ls -la cmd/ || echo "cmd directory not found"
RUN find . -name "*.go" | head -10

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o server ./cmd/server/main.go

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -S app -G app

WORKDIR /app

# Copy the binary from builder
COPY --from=builder --chown=app:app /app/server .
COPY --from=builder --chown=app:app /app/graph ./graph/

# Switch to non-root user
USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]