# Variables
APP_NAME = todo
BUILD_DIR = build
SRC_DIR = .
MAIN_SRC = $(SRC_DIR)/main.go
BINARY = $(BUILD_DIR)/$(APP_NAME)
DB_FILE = auth.db
DB_BACKUP_DIR = bak

# Helper function to backup database with timestamp
define backup_db
	@if [ -f "$(DB_FILE)" ]; then \
		TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
		NEW_NAME="$(1)/$${TIMESTAMP}-$(DB_FILE)"; \
		echo "Moving $(DB_FILE) to $${NEW_NAME}..."; \
		mv "$(DB_FILE)" "$${NEW_NAME}"; \
		echo "Database moved to $${NEW_NAME}"; \
	else \
		echo "Database file $(DB_FILE) not found"; \
	fi
endef

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

gencsrfkey:
	@if command -v openssl >/dev/null 2>&1; then \
		echo "CSRF Key: $$(openssl rand -base64 32)"; \
	elif command -v dd >/dev/null 2>&1; then \
		echo "CSRF Key: $$(dd if=/dev/urandom bs=32 count=1 2>/dev/null | base64)"; \
	else \
		echo "Neither openssl nor dd are available. Please install one of them."; \
		exit 1; \
	fi

# Set environment variables
# WIP: This is a workaround to be able to associate some styles to notifications and buttons but another approach will
# be used at the end.
setenv:
	@echo "Setting environment variables..."
	@export TODO_SERVER_WEB_HOST=localhost
	@export TODO_SERVER_WEB_PORT=8080
	@export TODO_SERVER_API_HOST=localhost
	@export TODO_SERVER_API_PORT=8081
	@export TODO_SERVER_INDEX_ENABLED=true
	@echo "Setting a CSRF key..."
	@export TODO_SEC_CSRF_KEY="NdZ7ULOe+NJ1bs5TzS51K+U4azOYQ6Wtv4CXlF6gJNM="
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
	@echo "Setting render errors..."
	@export TODO_RENDER_WEB_ERRORS="true"
	@export TODO_RENDER_API_ERRORS="true"
	@echo "Environment variables set."

# Clean the build directory
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Backup database in current directory
backup-db:
	$(call backup_db,.)

# Reset database by moving it to backup directory
reset-db:
	@mkdir -p $(DB_BACKUP_DIR)
	$(call backup_db,$(DB_BACKUP_DIR))
	@echo "A fresh database will be created on next application start"

# Phony targets
.PHONY: all build run runflags setenv clean backup-db reset-db