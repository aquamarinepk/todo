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

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@$(BINARY)

# Clean the build directory
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Phony targets
.PHONY: all build run clean
