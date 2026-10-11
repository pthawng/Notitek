# Notitek developer tasks. Run `make help` for the list.

SHELL := bash
.DEFAULT_GOAL := help

GOOSE_VERSION        ?= v3.28.0
GOLANGCI_LINT_VERSION ?= v2.14.0
COMPOSE              := docker compose -f deploy/compose.yaml
COMPOSE_DEV          := docker compose -f deploy/compose.dev.yaml --profile dev
# Hot reload rebuilds every running runtime on each change; run only what you work on.
SERVICES             ?= dev-api
POSTGRES_DSN         ?= postgres://notitek:notitek@localhost:5432/notitek?sslmode=disable
RUNTIMES             := api worker relay scheduler

GOOSE := go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION) -dir db/migrations postgres "$(POSTGRES_DSN)"

.PHONY: help
help: ## List targets
	@grep -hE '^[a-zA-Z0-9_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-16s %s\n", $$1, $$2}'

## ---- Code -------------------------------------------------------------------

.PHONY: fmt
fmt: ## Format code
	gofmt -w .

.PHONY: tidy
tidy: ## Tidy go.mod/go.sum
	go mod tidy

.PHONY: lint
lint: ## Run golangci-lint (includes module dependency rules)
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

.PHONY: test
test: ## Run unit tests
	go test ./...

.PHONY: test-race
test-race: ## Run unit tests with the race detector (needs cgo)
	go test -race ./...

.PHONY: test-integration
test-integration: ## Run integration tests against the compose stack
	NOTITEK_TEST_POSTGRES_DSN="$(POSTGRES_DSN)" go test -tags integration -count=1 ./...

.PHONY: contract
contract: ## Validate OpenAPI, examples and negative cases
	python -I scripts/validate_notification_contract.py

.PHONY: build
build: ## Build all runtimes into bin/
	@for r in $(RUNTIMES); do go build -trimpath -o bin/notitek-$$r ./cmd/$$r || exit 1; done

.PHONY: check
check: lint test contract ## Everything CI checks locally

## ---- Local stack --------------------------------------------------------------

.PHONY: up
up: ## Start Postgres, RabbitMQ and Redis
	$(COMPOSE) up -d --wait postgres rabbitmq redis

.PHONY: up-app
up-app: ## Start infrastructure and the four runtimes in containers
	$(COMPOSE) --profile app up -d --build

.PHONY: down
down: ## Stop the stack (keeps data)
	$(COMPOSE) --profile app down

.PHONY: reset
reset: ## Stop the stack and delete its data
	$(COMPOSE) --profile app down -v

.PHONY: logs
logs: ## Follow logs of the runtimes
	$(COMPOSE) --profile app logs -f $(RUNTIMES)

## ---- Hot reload (development) ----------------------------------------------

.PHONY: dev
dev: ## Hot reload (default dev-api): make dev SERVICES="dev-api dev-worker"
	$(COMPOSE_DEV) up -d --build $(SERVICES)
	@echo "Hot reload running. Logs: make dev-logs   Stop: make dev-down"

.PHONY: dev-logs
dev-logs: ## Follow logs of the hot-reload runtimes
	$(COMPOSE_DEV) logs -f $(SERVICES)

.PHONY: dev-down
dev-down: ## Stop the hot-reload stack (keeps data and Go caches)
	$(COMPOSE_DEV) down

## ---- Database -------------------------------------------------------------------

.PHONY: migrate-up
migrate-up: ## Apply migrations
	$(GOOSE) up

.PHONY: migrate-down
migrate-down: ## Roll back the last migration
	$(GOOSE) down

.PHONY: migrate-status
migrate-status: ## Show migration status
	$(GOOSE) status

.PHONY: migrate-new
migrate-new: ## Create a migration: make migrate-new name=add_outbox
	@test -n "$(name)" || (echo "usage: make migrate-new name=<snake_case>" && exit 1)
	$(GOOSE) -s create $(name) sql
