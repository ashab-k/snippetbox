FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy and download dependencies first (for better caching)
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
WORKDIR /app/cmd/web
RUN go build -o /app/snippetbox

# Use a smaller image for the final container
FROM alpine:latest

# Install CA certificates for any HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/snippetbox /app/snippetbox

# Copy UI templates and static files
COPY --from=builder /app/ui /app/ui

# Expose the application port
EXPOSE 4000

# Run the application
CMD ["/app/snippetbox"]