.DEFAULT_GOAL := help

.PHONY: help
help: ## Print this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

GOFLAGS := -buildvcs=false

.PHONY: build
build: ## Build the project
	go build $(GOFLAGS) -o bin/mdcli .

.PHONY: install
install: build ## Install the binary to ~/.local/bin
	mkdir -p ~/.local/bin
	ln -sf $(CURDIR)/bin/mdcli ~/.local/bin/mdcli

.PHONY: test-build
test-build: ## Build test code
	go test $(GOFLAGS) -c ./...

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

.PHONY: fmt lint
fmt: ## Format go code
	GOFLAGS="$(GOFLAGS)" golangci-lint run
	@files="$$(gofmt -l -w .)"; \
	if [ -n "$$files" ]; then \
		echo "Reformatted these files:"; \
		echo "$$files"; \
		exit 1; \
	fi
	go vet ./...

lint: fmt # Alias for fmt

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
