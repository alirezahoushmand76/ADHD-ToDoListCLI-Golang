# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /todolist ./cmd/todolist

# Final stage
FROM alpine:latest

# Install bash for better shell experience
RUN apk --no-cache add bash

# Create a non-root user
RUN adduser -D -h /home/appuser appuser
USER appuser
WORKDIR /home/appuser

# Copy the binary from the builder stage
COPY --from=builder /todolist /usr/local/bin/todolist

# Create data directory
RUN mkdir -p /home/appuser/.todolist/backups

# Set entrypoint
ENTRYPOINT ["todolist"]

# Default command
CMD ["--help"] 