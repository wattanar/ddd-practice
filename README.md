# Task Manager API

A task management REST API built with **Clean Architecture + DDD + Hexagonal** in Go.

## Architecture

```
Infrastructure → Application → Domain
   (adapters)     (use cases)    (core)
```

- **Domain** — Pure Go, zero external dependencies. Task entity, value objects, domain events, repository port.
- **Application** — Use cases (Create, Get, List, Update, Delete). Depends on domain only.
- **Infrastructure** — PostgreSQL adapter, HTTP handlers (generated from OpenAPI). Implements ports.

## Quickstart

```bash
# Start PostgreSQL
make dc-up

# Run migrations (Liquibase via Docker)
make migrate

# Run server on :8080
make run

# Run all 45 tests
make test
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST`   | `/api/v1/tasks`      | Create a task |
| `GET`    | `/api/v1/tasks`      | List tasks (optional `?status=` & `?priority=` filters) |
| `GET`    | `/api/v1/tasks/{id}` | Get a task by ID |
| `PATCH`  | `/api/v1/tasks/{id}` | Update a task |
| `DELETE` | `/api/v1/tasks/{id}` | Delete a task |

OpenAPI spec at [`api/openapi.yaml`](api/openapi.yaml).

## Project Structure

```
├── api/
│   ├── openapi.yaml              # OpenAPI 3.0.3 spec
│   └── spec/server.gen.go        # Generated types + ServerInterface
├── cmd/server/main.go            # Composition root / entry point
├── internal/
│   ├── domain/                   # Business logic (zero external deps)
│   │   ├── task/                 # Task aggregate: entity, value objects, events, repository port
│   │   └── errors/              # Domain error types
│   ├── application/              # Use cases (depends on domain only)
│   │   └── task/                 # Create, Get, List, Update, Delete use cases
│   └── infrastructure/          # Adapters (implements ports)
│       ├── persistence/postgres/ # PostgreSQL adapter using pgx v5
│       └── http/                 # HTTP handlers (implements ServerInterface)
├── migrations/                   # Liquibase changelog + SQL
├── docker-compose.yml            # PostgreSQL 17
└── Makefile                      # Build, test, codegen, migrate targets
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `run` | Run server |
| `build` | Build binary |
| `generate` | Regenerate OpenAPI code + mocks |
| `generate/openapi` | Generate server/types from OpenAPI spec |
| `generate/mocks` | Generate mocks with mockery |
| `test` | Run all tests |
| `test/unit` | Domain + application tests only |
| `test/http` | HTTP handler tests |
| `test/integration` | PostgreSQL integration tests (`-tags=integration`) |
| `migrate` | Apply Liquibase migrations |
| `migrate/rollback` | Rollback last changeset |
| `dc-up` | Start PostgreSQL via Docker Compose |
| `dc-down` | Stop containers |
| `lint` | Run staticcheck |
| `clean` | Remove generated + build artifacts |

## Testing Layers

| Layer | Count | Approach |
|-------|-------|----------|
| Domain | 11 | Pure unit tests, no deps |
| Application | 8 | Mockery-generated mocks |
| HTTP | 6 | httptest with mocked use cases |
| Integration | 1 | Real PostgreSQL (`DATABASE_URL` env) |

```bash
make test        # 45 tests, all layers
make test/unit   # 19 tests, fast (< 1s)
```

## Configuration

| Env Var | Default | Description |
|---------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DATABASE_URL` | `postgres://taskmanager:taskmanager@localhost:5432/taskmanager?sslmode=disable` | PostgreSQL connection string |

## Dependencies

- **HTTP** — Go 1.22+ stdlib (`net/http`, enhanced ServeMux with method routing)
- **Database** — `github.com/jackc/pgx/v5`
- **Codegen** — `oapi-codegen` (OpenAPI), `mockery` (mocks)
- **Migrations** — Liquibase via Docker

## Domain Model

| Value Object | Rules |
|---|---|
| `TaskID` | UUID v4 |
| `Title` | 1–200 characters |
| `Description` | 0–2000 characters |
| `TaskStatus` | `todo` → `in_progress` → `done` (reopenable) |
| `Priority` | `low`, `medium`, `high`, `critical` |

Domain events (`task.created`, `task.status_changed`, `task.title_changed`, `task.deleted`) are collected on the entity and available via `PullEvents()`.
