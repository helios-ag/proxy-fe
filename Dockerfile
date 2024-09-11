# Build stage
FROM golang:1.23.0-bookworm AS builder
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code and config
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM builder
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 3000

# Command to run
CMD ["./main"]