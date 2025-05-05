# ---------- Build Stage ----------
FROM golang:1.24 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy dependency files and download modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with a specific name
RUN CGO_ENABLED=0 GOOS=linux go build -o advertising ./cmd

# ---------- Run Stage ----------
FROM alpine:latest

# Set working directory in the final container
WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /app/advertising .

# Expose the application's port
EXPOSE 8080

# Start the application
CMD ["./advertising"]
