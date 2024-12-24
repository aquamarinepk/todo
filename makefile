# Variables
APP_NAME = todo
BUILD_DIR = build
SRC_DIR = .
MAIN_SRC = $(SRC_DIR)/main.go
BINARY = $(BUILD_DIR)/$(APP_NAME)

# Default target
all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY) $(MAIN_SRC)
	@echo "Build complete: $(BINARY)"

# Run the application with environment variables
run: setenv build
	@echo "Running $(APP_NAME) with environment variables..."
	@$(BINARY)

# Run the application with command-line flags
runflags: build
	@echo "Running $(APP_NAME) with command-line flags..."
	@$(BINARY) -server.web.host=localhost -server.web.port=9080 -server.api.host=localhost -server.api.port=9081

# Set environment variables
setenv:
	@echo "Setting environment variables..."
	@export TODO_SERVER_WEB_HOST=localhost
	@export TODO_SERVER_WEB_PORT=8080
	@export TODO_SERVER_API_HOST=localhost
	@export TODO_SERVER_API_PORT=8081
	@echo "Environment variables set."

# Clean the build directory
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Phony targets
.PHONY: all build run run_flags setenv clean