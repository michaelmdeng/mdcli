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

.PHONY: fmt lint
fmt: ## Format go code
	golangci-lint run
	go vet ./...
	go fmt ./...

lint: ## Alias for fmt
	golangci-lint run
	go vet ./...
	go fmt ./...

.PHONY: test
test: ## Run tests
	go test ./...

TAG ?= latest
REGISTRY_TAG ?= ghcr.io/michaelmdeng/mdcli/mdcli:$(TAG)

.PHONY: build-image
build-image: ## Build the docker image, accepts optional tag. Usage: make build-image tag=v0.0.1
	docker build -t $(REGISTRY_TAG) .

.PHONY: run-image
run-image: ## Run the docker image, accepts optional tag. Usage: make run-image tag=v0.0.1
	docker run -it $(REGISTRY_TAG)

.PHONY: publish-image
publish-image: ## Publish the docker image to ghcr.io, requires tag. Usage: make publish-image tag=v0.0.1
	docker push $(REGISTRY_TAG)
