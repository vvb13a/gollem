# App developer tasks. Run `make` for the list of targets.
#
# The Vue app lives in ./frontend; npm's --prefix flag runs its scripts from
# this directory, so no `cd` is needed.

.DEFAULT_GOAL := help

# Capture this file before -include adds .env to MAKEFILE_LIST (the help target
# greps only this file).
THIS_MAKEFILE := $(firstword $(MAKEFILE_LIST))

# Optional local overrides.
-include .env

DEV_PORT ?= 5173
APP_ADDR ?= :8080
APP_NAME ?= gollem
APP_DATA_DIR ?= $(CURDIR)/.dev/data
GOARCH ?= amd64
VERSION ?= dev

NPM := npm --prefix frontend

.PHONY: help deps install web dev api build-web build release test tidy clean

help: ## List available targets
	@echo "App targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(THIS_MAKEFILE) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

deps:
	@test -d frontend/node_modules || $(MAKE) install

install: ## Install frontend dependencies
	$(NPM) ci

web: deps ## Run the Vite dev server only
	$(NPM) run dev -- --port $(DEV_PORT)

dev: deps ## Run the Go API and Vite together
	@mkdir -p "$(APP_DATA_DIR)"
	@echo "API $(APP_ADDR), Vite http://localhost:$(DEV_PORT), data $(APP_DATA_DIR)"
	@trap 'kill 0' INT TERM; \
		APP_ADDR="$(APP_ADDR)" \
		APP_NAME="$(APP_NAME)" \
		APP_DATA_DIR="$(APP_DATA_DIR)" \
		go run . & \
		$(NPM) run dev -- --port $(DEV_PORT) & \
		wait

api: ## Run the Go server only
	@mkdir -p "$(APP_DATA_DIR)"
	APP_ADDR="$(APP_ADDR)" APP_NAME="$(APP_NAME)" APP_DATA_DIR="$(APP_DATA_DIR)" go run .

build-web: deps ## Build the Vue app into frontend/dist
	@find frontend/dist -mindepth 1 ! -name .gitkeep -exec rm -rf {} +
	$(NPM) run build

build: build-web ## Build the Go binary into bin/app
	@mkdir -p bin
	go build -ldflags="-X main.version=$(VERSION)" -o bin/app .

release: build-web ## Cross-compile a static linux binary into bin/app-linux-$(GOARCH)
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" \
		-o bin/app-linux-$(GOARCH) .

test: ## Run the Go tests
	go test ./...

tidy: ## Tidy Go modules
	go mod tidy

clean: ## Remove build artifacts (keeps frontend/dist/.gitkeep)
	rm -rf bin
	find frontend/dist -mindepth 1 ! -name .gitkeep -exec rm -rf {} +
