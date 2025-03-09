# syntax=docker/dockerfile:1

# --- Builder Stage ---
FROM golang:1.18 AS builder

WORKDIR /src

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build the Go application. Adjust the binary name ("app") as needed.
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# --- Final Stage ---
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /src/app .

# Expose port 8080 (or change as needed)
EXPOSE 8080

# Run the application
ENTRYPOINT ["./app"]
