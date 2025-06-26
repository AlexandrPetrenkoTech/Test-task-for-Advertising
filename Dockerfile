FROM golang:1.24 AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd

FROM alpine:latest

WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/main .

# Copy configuration files
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/.env .

# Copy generated Swagger documentation
COPY --from=builder /app/docs ./docs

# Expose port (optional, for documentation)
# EXPOSE 8080

# Run the application
CMD ["./main"]
