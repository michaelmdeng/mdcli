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


.PHONY: aider
aider:
	uv tool run --from aider-chat aider --conf .aider.conf.yaml $(filter-out $@,$(MAKECMDGOALS))
