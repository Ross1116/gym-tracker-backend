# Variables
BINARY_NAME=bin/gym-tracker-backend

# Build the binary
build:
	@go build -o $(BINARY_NAME)

# Run the built binary
run: build
	@./$(BINARY_NAME)

# Run tests with verbose output
test:
	@go test -v ./...

# Clean build artifacts
clean:
	@rm -rf $(BINARY_NAME)
