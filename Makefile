# Makefile for Advertising Project

# Name of the final binary
BINARY_NAME=advertising

# Default target
.PHONY: all
all: build

# Run the linter
.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

# Build the Go binary
.PHONY: build
build:
	@echo "Building Go project..."
	CGO_ENABLED=0 GOOS=linux go build -o $(BINARY_NAME) ./cmd

# Run the application locally
.PHONY: run
run:
	@echo "Running app..."
	go run ./cmd

# Run unit tests
.PHONY: test
test:
	@echo "Running unit tests..."
	go test ./... -v

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t advertising-app .

# Run Docker container
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 advertising-app

# Clean up built files
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
