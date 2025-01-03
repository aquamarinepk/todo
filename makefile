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
# WIP: This is a workaround to be able to assoiciate some styles to notifications and buttons but probably another approach will be used at the end.
setenv:
	@echo "Setting environment variables..."
	@export TODO_SERVER_WEB_HOST=localhost
	@export TODO_SERVER_WEB_PORT=8080
	@export TODO_SERVER_API_HOST=localhost
	@export TODO_SERVER_API_PORT=8081
	@echo "Setting notification styles..."
	@export TODO_NOTIFICATION_SUCCESS_STYLE="bg-green-600 text-white px-4 py-2 rounded"
	@export TODO_NOTIFICATION_INFO_STYLE="bg-blue-600 text-white px-4 py-2 rounded"
	@export TODO_NOTIFICATION_WARN_STYLE="bg-yellow-600 text-white px-4 py-2 rounded"
	@export TODO_NOTIFICATION_ERROR_STYLE="bg-red-600 text-white px-4 py-2 rounded"
	@export TODO_NOTIFICATION_DEBUG_STYLE="bg-gray-600 text-white px-4 py-2 rounded"
	@echo "Setting button styles..."
	@export TODO_BUTTON_STYLE_STANDARD="bg-gray-600 text-white px-4 py-2 rounded"
	@export TODO_BUTTON_STYLE_BLUE="bg-blue-600 text-white px-4 py-2 rounded"
	@export TODO_BUTTON_STYLE_RED="bg-red-600 text-white px-4 py-2 rounded"
	@export TODO_BUTTON_STYLE_GREEN="bg-green-600 text-white px-4 py-2 rounded"
	@export TODO_BUTTON_STYLE_YELLOW="bg-yellow-600 text-white px-4 py-2 rounded"
	@echo "Environment variables set."

# Clean the build directory
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Phony targets
.PHONY: all build run run_flags setenv clean
