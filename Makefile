build: ### build
	go build -v ./cmd/app
	go build -v ./cmd/client
.PHONY: build

test: ### run test
	go test -v -cover -race -timeout 30s ./internal/...
.PHONY: test

fmt: ### run format
	go fmt ./...
.PHONY: fmt

vet: ### run linter
	go vet ./...
.PHONY: vet 

mock: ### run mockgen
	mockgen -source=./internal/usecase/interfaces.go -destination=./internal/usecase/mocks.go -package=usecase
.PHONY: mock 

migrate-create: ### create new migration
	migrate create -ext sql -dir migrations 'migrate_name'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(DATABASE_DSN)' up
.PHONY: migrate-up

coverage: ### run html cover file
	go test -race -timeout 30s -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out
.PHONY: coverage

.DEFAULT_GOAL := build
