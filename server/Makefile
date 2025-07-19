# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

# Build configuration
BINARY_NAME := proompt
BUILD_DIR := build
CMD_DIR := cmd/proompt

.PHONY: build test test-unit test-integration clean dev help

help: ## Show this help message
	@echo "$(BLUE)Available targets:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2}'

build: ## Build binary to build/proompt
	@echo "$(BLUE)Building $(BINARY_NAME)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@go build -buildvcs=false -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(RESET)"

test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests only
	@echo "$(BLUE)Running unit tests...$(RESET)"
	@go test -v ./internal/...
	@echo "$(GREEN)Unit tests complete$(RESET)"

test-integration: ## Run integration tests
	@echo "$(YELLOW)Integration tests not implemented yet$(RESET)"

clean: ## Clean build directory
	@echo "$(BLUE)Cleaning build directory...$(RESET)"
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)Clean complete$(RESET)"

dev: build ## Build and run in development mode
	@echo "$(BLUE)Starting development server...$(RESET)"
	@./$(BUILD_DIR)/$(BINARY_NAME)

.DEFAULT_GOAL := help