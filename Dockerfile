# Use the official Golang image as the base
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Copy .env file into the container
COPY .env /app/.env

# Build the Go application
RUN go build -o main ./cmd/main.go

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
