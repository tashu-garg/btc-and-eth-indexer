## Multi-stage Dockerfile for Go Indexer Backend

# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install git (required for go modules that need it)
RUN apk add --no-cache git

# Cache go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . ./

# Build the binary
# Adjust the output path/name if needed
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/indexer cmd/server/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.20
# Add ca-certificates for HTTPS calls
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser:appgroup

WORKDIR /app
COPY --from=builder /app/bin/indexer ./indexer
# Optionally copy .env.example or other needed files
# COPY --from=builder /app/.env.example ./

# Expose the port (default from config, typically 8080)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./indexer"]

# Optional: define default command arguments (none needed)
# CMD []
