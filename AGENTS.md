# Task Manager — Agent Guide

## Architecture (Dependency Flow)

```
infrastructure/http → application/task → domain/task
                                       ↗
infrastructure/persistence/postgres ──┘
```

- **`domain/task`** — zero external deps (no pgx, no HTTP). All business rules here.
- **`application/task`** — depends only on `domain/task`. One use case per file.
- **`infrastructure/persistence/postgres`** — implements `task.TaskRepository` port. Returns `nil, nil` on not-found (not an error).
- **`infrastructure/http`** — implements generated `ServerInterface`. Translates OpenAPI DTOs ↔ domain objects. Must NOT leak domain types to the API spec or vice versa.

## Generated Code (Always Regenerate)

| Source | Generator | Output |
|--------|-----------|--------|
| `api/openapi.yaml` | `oapi-codegen -generate types,std-http` | `api/spec/server.gen.go` |
| `internal/domain/task/repository.go` | `mockery --config .mockery.yaml` | `internal/domain/task/mocks/` |

`make generate` runs both. The `test` and `run` targets depend on `generate`. Run `make generate` before committing after changing any interface or the OpenAPI spec.

## Key Commands

```bash
make test            # generate + go test -v -count=1 ./...
make test/unit       # skip infra: domain + application only
make test/integration  # -tags=integration, requires DATABASE_URL + postgres
make run             # generate + go run ./cmd/server
make lint            # generate + staticcheck ./...
make clean           # rm -rf bin/ api/spec/ mocks/
```

## Testing Rules

- **Integration tests** have build tag `//go:build integration` — won't run without `-tags=integration`
- **Application tests** use mockery-generated mocks in `internal/domain/task/mocks/`
- **HTTP tests** use `httptest.Server` with `spec.HandlerFromMux()` — test through the generated router
- **Repository interface returns `nil, nil`** on not-found. Use cases must check for nil after `FindByID`.
- Use `count=1` to disable test caching (`make test` does this).

## Domain Conventions

- **Value objects** validate on construction via `NewXxx()` — return `error` for invalid state. No setters.
- **Entity behavior** mutates state and emits events via `PullEvents()`. `ReconstituteTask()` is for DB reads only — does NOT emit events.
- **Status transitions**: `todo → in_progress → done` (reopenable: `done → in_progress`). Invalid transitions return `InvalidTransition` error.
- **`FindByID` returning nil** means not-found. Always check: `if t == nil { return nil, &NotFound{...} }`.

## Composition Root

`cmd/server/main.go` is the only place where layers are wired:
```go
pool := pgxpool.New(ctx, dsn)
repo := postgres.NewTaskRepository(pool)
uc := taskapp.NewCreateTaskUseCase(repo)  // repeat per use case
handler := taskhttp.NewTaskHandler(createUC, getUC, ...)
mux := spec.HandlerFromMux(handler, http.NewServeMux())
```

Config from env vars (`PORT`, `DATABASE_URL`). No config file. Hardcoded defaults for local dev.

## Go Tooling

- **Requires Go 1.22+** (uses enhanced `http.ServeMux` with `{id}` path params and method routing). Go 1.26.3 in go.mod.
- **No DI framework** — manual wiring in main.
- Lint: `staticcheck ./...`
- Only test deps beyond std: `pgx/v5`, `uuid`, `oapi-codegen/runtime`, `testify`.

## Adding a Feature

1. Add/change OpenAPI spec → `make generate/openapi`
2. Add domain logic (entity behavior, value objects) → domain package
3. Add/update repository interface → domain package → `make generate/mocks`
4. Add use case → application package
5. Add/update handler → infrastructure/http → implements `ServerInterface`
6. Wire in `cmd/server/main.go`
7. Tests: domain (unit), application (mock), HTTP (httptest), integration (tagged)
