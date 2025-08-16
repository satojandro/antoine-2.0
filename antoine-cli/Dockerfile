FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh antoine

# Set working directory
WORKDIR /home/antoine

# Copy binary from builder
COPY --from=builder /app/bin/antoine /usr/local/bin/antoine

# Copy default config
COPY configs/default.yaml /home/antoine/.antoine.yaml

# Change ownership
RUN chown -R antoine:antoine /home/antoine

# Switch to non-root user
USER antoine

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD antoine config show || exit 1

# Default command
ENTRYPOINT ["antoine"]
CMD ["--help"]
