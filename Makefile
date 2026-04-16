BINARY_SERVER   := bin/server
BINARY_QUICKGEN := bin/quickgen
MODULE          := github.com/samaita/quick-go

.PHONY: all build build-server build-quickgen run dev tidy clean help

all: build

## build: Build both server and quickgen binaries into bin/
build: build-server build-quickgen

build-server:
	@mkdir -p bin
	go build -o $(BINARY_SERVER) ./cmd/server

build-quickgen:
	@mkdir -p bin
	go build -o $(BINARY_QUICKGEN) ./cmd/quickgen

## run: Run the server (requires .env)
run: build-server
	./$(BINARY_SERVER)

## dev: Run the server with live .env reload using go run
dev:
	go run ./cmd/server

## gen: Run quickgen against a schema file. Usage: make gen SCHEMA=path/to/schema.sql
gen: build-quickgen
	./$(BINARY_QUICKGEN) --schema $(SCHEMA) --out .

## tidy: Tidy and verify Go modules
tidy:
	go mod tidy
	go mod verify

## clean: Remove built binaries
clean:
	rm -rf bin/

## help: Show this help message
help:
	@grep -E '^##' Makefile | sed 's/## //'
