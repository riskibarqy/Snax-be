FROM golang:1.23.6-alpine AS builder

WORKDIR /app

# Copy go mod files from root directory
COPY ../go.mod ../go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY .. .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/url-shortener /app/url-shortener/cmd/main.go && \
    chmod 755 /usr/local/bin/url-shortener && \
    ls -la /usr/local/bin/url-shortener

# Create final image
FROM alpine:latest

# Install necessary tools
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and migrations
COPY --from=builder /usr/local/bin/url-shortener /usr/local/bin/url-shortener
COPY --from=builder /app/sqlc/schema /app/sqlc/schema

# Verify binary permissions
RUN ls -la /usr/local/bin/url-shortener && \
    chmod 755 /usr/local/bin/url-shortener

# Set environment variables
ENV TZ=UTC \
    SERVICE_NAME=url-shortener \
    ENV=production \
    PORT=8080

# Expose port
EXPOSE 8080

# Run the application
CMD ["/usr/local/bin/url-shortener"] 