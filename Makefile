# Opgecanceld Blocklist — YouTube ad domain blocklist and discovery tool
# Follows Standard Go Project Layout: https://github.com/golang-standards/project-layout

SHELL := /bin/bash
.DEFAULT_GOAL := help

# Build
BINARY := discover
MAIN := ./cmd/discover
BIN_DIR := bin
SCRIPTS := scripts
OUT := $(BIN_DIR)/$(BINARY)

.PHONY: help build run dev test lint fmt clean install

help:
	@echo "Opgecanceld Blocklist — YouTube ad domain discovery"
	@echo ""
	@echo "Targets:"
	@echo "  build   build binary to $(OUT)"
	@echo "  run     build and run the binary"
	@echo "  dev     clean, then build and run"
	@echo "  test    run tests"
	@echo "  lint    run golangci-lint"
	@echo "  fmt     format Go code and tidy modules"
	@echo "  clean   remove built binary"
	@echo "  install build and install to $$(go env GOPATH)/bin"
	@echo ""

build:
	@mkdir -p $(BIN_DIR)
	BINARY=$(OUT) MAIN=$(MAIN) $(SCRIPTS)/build.sh

run: build
	$(OUT)

dev: clean run

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w .
	go mod tidy

clean:
	rm -f $(OUT)
	@rmdir $(BIN_DIR) 2>/dev/null || true

install: build
	go install $(MAIN)
