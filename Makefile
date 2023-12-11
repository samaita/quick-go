# Makefile for building and running Go project

# Binary name
BINARY_NAME := quickgo

# Directories
CMD_DIR := ./cmd
TMP_DIR := ./tmp

# Build command
build:
	@mkdir -p $(TMP_DIR)
	@cd $(CMD_DIR) && go build -o ../$(TMP_DIR)/$(BINARY_NAME)

# Run command
run: build
	@cd $(TMP_DIR) && ./$(BINARY_NAME)

# Clean command
clean:
	@rm -rf $(TMP_DIR)

.PHONY: build run clean