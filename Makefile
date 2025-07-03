.DEFAULT_GOAL := help

.PHONY: help
help: ## Print this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the project
	go build -o bin ./...

.PHONY: test-build
test-build: ## Build test code
	go test -c ./...

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

.PHONY: fmt
fmt: ## Format go code
	go fmt ./...

.PHONY: test
test: ## Run tests
	go test ./...
