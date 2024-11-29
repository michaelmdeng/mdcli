.PHONY: build
build:
	go build -o bin/mdcli ./...

.PHONY: tidy
tidy:
	base=$$(echo $$PWD); \
	go mod tidy; \
	for f in $$(ls src); do cd $$base/src/$$f && go mod tidy; done;

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./... src/...
