.PHONY: build test clean install

APP_NAME = hyperion
MAIN_FILE = main.go

# Go related variables.
GOBIN = $(GOPATH)/bin

# Build the Go app
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(APP_NAME) $(MAIN_FILE)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f $(APP_NAME)
	@go clean

# Install the application
install: build
	@echo "Installing $(APP_NAME)..."
	@mv $(APP_NAME) $(GOBIN)/$(APP_NAME)
	@echo "Installation complete. You can now run '$(APP_NAME)' from anywhere."

# Run the application with default settings
run: build
	@echo "Running $(APP_NAME)..."
	@./$(APP_NAME)

# Run the application with all features enabled
run-full: build
	@echo "Running $(APP_NAME) with all features enabled..."
	@./$(APP_NAME) --show-files --unicode --color --show-stats --stat-table --chart

# Help
help:
	@echo "Available commands:"
	@echo "  make build     - Build the application"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make install   - Install the application to GOBIN"
	@echo "  make run       - Build and run with default settings"
	@echo "  make run-full  - Build and run with all features enabled"