# Use the Go 1.23.2 version
FROM golang:1.23.2-alpine AS builder

# Set a flexible working directory
WORKDIR /app/go-startup

# Copy Go module files and download dependencies# Copy the entire application
COPY . .
RUN go mod tidy


# Build the Go application
RUN go build -o main .

# Use a minimal image to serve the Go app
FROM alpine:latest

WORKDIR /app/go-startup

# Copy the compiled binary
COPY --from=builder /app/go-startup/main .
COPY --from=builder /app/go-startup/.env .

# Copy everything from the web and images directories
COPY --from=builder /app/go-startup/web/ /app/go-startup/web/
COPY --from=builder /app/go-startup/images/ /app/go-startup/images/

# Expose the API port
EXPOSE 80

# Start the Go application
CMD ["./main"]
