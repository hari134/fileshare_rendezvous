# Stage 1: Build the Go binary
FROM golang:1.21.0 AS builder

WORKDIR /app

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Stage 2: Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy only the built binary from the previous stage
COPY --from=builder /app/app .

# Expose the port that the application will run on
EXPOSE 9000

# Command to run the executable
CMD ["./app"]