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

.PHONY: help build build-filters run dev test lint lint-yaml fmt format format-check prettier renovate clean install

help:
	@echo "Opgecanceld Blocklist — YouTube ad domain discovery"
	@echo ""
	@echo "Targets:"
	@echo "  build         build binary to $(OUT)"
	@echo "  build-filters generate AdGuard/uBlock filter list from blocklist"
	@echo "  run          build and run the binary"
	@echo "  dev          clean, then build and run"
	@echo "  test         run tests"
	@echo "  lint         run golangci-lint"
	@echo "  lint-yaml    run yamllint on YAML files"
	@echo "  fmt          format Go code and tidy modules"
	@echo "  format       format all files with Prettier (same version as CI)"
	@echo "  format-check check formatting without making changes"
	@echo "  prettier     alias for format"
	@echo "  renovate     run renovate (dry-run)"
	@echo "  clean        remove built binary"
	@echo "  install      build and install to $$(go env GOPATH)/bin"
	@echo ""

build:
	@mkdir -p $(BIN_DIR)
	BINARY=$(OUT) MAIN=$(MAIN) $(SCRIPTS)/build.sh

build-filters:
	python3 $(SCRIPTS)/build-filters.py

run: build
	$(OUT)

dev: clean run

test:
	go test ./...

lint:
	golangci-lint run

lint-yaml:
	@echo "Running yamllint..."
	@command -v yamllint >/dev/null 2>&1 || { echo "yamllint not found. Install with: brew install yamllint or pip install yamllint"; exit 1; }
	@yamllint .
	@echo "yamllint done."

fmt:
	gofmt -s -w .
	go mod tidy

# Prettier version aligned with CI
PRETTIER := npx --yes prettier@3.3.2

format:
	@echo "Formatting files with Prettier..."
	@$(PRETTIER) --write .

format-check:
	@echo "Checking file formatting..."
	@$(PRETTIER) --check .

prettier: format

renovate:
	npx -y renovate --dry-run

clean:
	rm -f $(OUT)
	@rmdir $(BIN_DIR) 2>/dev/null || true

install: build
	go install $(MAIN)
