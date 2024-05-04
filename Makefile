up: ## run docker compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: up

down: ## down docker compose
	docker-compose down --remove-orphans
.PHONY: down

build: ### build
	mkdir -p bin
	go build -v -o ./bin ./cmd/app 
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

migrate-create: ### create new migration
	migrate create -ext sql -dir migrations 'migrate_name'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(DATABASE_DSN)' up
.PHONY: migrate-up

coverage: ### run html cover file
	go test -race -timeout 30s -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out
	rm coverage.out
.PHONY: coverage

.DEFAULT_GOAL := build
