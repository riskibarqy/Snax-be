FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./url-shortener/cmd/main.go

# Create final image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/url-shortener .

# Copy migrations
COPY --from=builder /app/sqlc/schema ./sqlc/schema

# Install necessary tools
RUN apk add --no-cache ca-certificates tzdata

# Set environment variables
ENV TZ=UTC \
    SERVICE_NAME=url-shortener

# Expose port
EXPOSE 8080

# Run the application
CMD ["./url-shortener"] 