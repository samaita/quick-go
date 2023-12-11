# Makefile for building and running Go project

include .env
export

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

migrate-up:
	@migrate -path migrations -database "postgres://$(DATABASE.USER):$(DATABASE.PASSWORD)@$(DATABASE.HOST):$(DATABASE.PORT)/$(DATABASE.DBNAME)?sslmode=$(DATABASE.SSLMODE)" up

migrate-down:
	@migrate -path migrations -database "postgres://$(DATABASE.USER):$(DATABASE.PASSWORD)@$(DATABASE.HOST):$(DATABASE.PORT)/$(DATABASE.DBNAME)?sslmode=$(DATABASE.SSLMODE)" down

migrate-force:
	$(eval MIGRATE_VERSION := $(shell read -p "Enter the version to force: " && echo "$$REPLY"))
	@migrate -path migrations -database "postgres://$(DATABASE.USER):$(DATABASE.PASSWORD)@$(DATABASE.HOST):$(DATABASE.PORT)/$(DATABASE.DBNAME)?sslmode=$(DATABASE.SSLMODE)" force $(MIGRATE_VERSION)

migrate-version:
	@migrate -path migrations -database "postgres://$(DATABASE.USER):$(DATABASE.PASSWORD)@$(DATABASE.HOST):$(DATABASE.PORT)/$(DATABASE.DBNAME)?sslmode=$(DATABASE.SSLMODE)" version

migrate-create:
	@read -p "Enter the name of the migration file: " migration_name; \
	migrate create -ext sql -dir migrations -seq $$migration_name

.PHONY: build run clean