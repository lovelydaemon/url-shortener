.PHONY: build
build:
	go build -v ./cmd/shortener

.PHONY: test
test:
	go test -v -cover -race -timeout 30s ./internal/...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...


.DEFAULT_GOAL := build
