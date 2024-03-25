.PHONY: build
build:
	go build -v ./cmd/shortener

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...


.DEFAULT_GOAL := build
