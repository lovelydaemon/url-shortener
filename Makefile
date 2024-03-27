.PHONY: build
build:
	go build -v ./cmd/app

.PHONY: test
test:
	go test -v -cover -race -timeout 30s ./internal/...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: coverage
coverage:
	go test -race -timeout 30s -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out

.DEFAULT_GOAL := build
