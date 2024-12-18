.PHONY: build
build:
	go build -o bin ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...
