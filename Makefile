APP_NAME        := taskmanager
BUILD_DIR       := ./bin
MIGRATIONS_DIR  := ./migrations
LIQUIBASE_IMAGE := liquibase/liquibase:latest
LIQUIBASE_CMD   := docker run --rm -v "$(PWD)/$(MIGRATIONS_DIR):/liquibase/changelog" \
                    $(LIQUIBASE_IMAGE) \
                    --defaultsFile=/liquibase/changelog/liquibase.properties

.PHONY: help run build test test/unit test/http test/integration \
        generate generate/openapi generate/mocks migrate migrate/rollback migrate/status \
        dc-up dc-down lint clean

help: ## print targets
	@grep -E '^[a-zA-Z/_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## === Build & Run ===

run: generate ## run server
	go run ./cmd/server

build: generate ## build binary
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

## === Code Generation ===

generate: generate/openapi generate/mocks ## run all code generators

generate/openapi: ## generate OpenAPI server + types
	oapi-codegen -generate types,std-http -package spec -o api/spec/server.gen.go api/openapi.yaml

generate/mocks: ## generate mocks with mockery
	mockery --config .mockery.yaml

## === Tests ===

test: generate ## run all tests
	go test -v -count=1 ./...

test/unit: ## domain + application tests only
	go test -v -count=1 ./internal/domain/... ./internal/application/...

test/http: ## HTTP handler tests
	go test -v -count=1 ./internal/infrastructure/http/...

test/integration: ## postgres integration tests
	go test -v -count=1 -tags=integration ./internal/infrastructure/persistence/...

## === Database (Liquibase) ===

migrate: ## apply pending migrations
	$(LIQUIBASE_CMD) update

migrate/rollback: ## rollback last changeset
	$(LIQUIBASE_CMD) rollback-count 1

migrate/status: ## show migration status
	$(LIQUIBASE_CMD) status

## === Docker ===

dc-up: ## start PostgreSQL
	docker compose up -d postgres

dc-down: ## stop all containers
	docker compose down

## === Lint & Clean ===

lint: generate ## run linter
	staticcheck ./...

clean: ## remove generated + build artifacts
	rm -rf $(BUILD_DIR) api/spec internal/domain/task/mocks
