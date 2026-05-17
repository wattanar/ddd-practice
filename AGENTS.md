# Task Manager — Agent Guide

## Architecture (Dependency Flow)

```
infrastructure/http → application/{task,department} → domain/{task,department}
                                                     ↗
infrastructure/persistence/postgres ──────────────────┘
```

Two bounded contexts, same layered pattern:
- **`domain/{task,department}`** — zero external deps. All business rules here.
- **`application/{task,department}`** — depends only on its domain package. One use case per file.
- **`infrastructure/persistence/postgres`** — implements repository ports. Returns `nil, nil` on not-found (not an error).
- **`infrastructure/http`** — implements generated `ServerInterface`. Two handlers (`TaskHandler`, `DepartmentHandler`) composed via `APIHandler`. Translates OpenAPI DTOs ↔ domain objects.

## Generated Code (Always Regenerate)

| Source | Generator | Output |
|--------|-----------|--------|
| `api/openapi.yaml` | `oapi-codegen -generate types,std-http` | `api/spec/server.gen.go` |
| `internal/domain/task/repository.go` | `mockery --config .mockery.yaml` | `internal/domain/task/mocks/` |
| `internal/domain/department/repository.go` | `mockery --config .mockery.yaml` | `internal/domain/department/mocks/` |

`make generate` runs all generators. The `test`, `build`, and `run` targets depend on `generate`. Regenerate after changing any interface or the OpenAPI spec.

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
- **Application tests** use mockery-generated mocks in `internal/domain/task/mocks/` and `internal/domain/department/mocks/`
- **HTTP tests** use `httptest.Server` with `spec.HandlerFromMux()` — test through the generated router
- **Repository interface returns `nil, nil`** on not-found. Use cases must check for nil after `FindByID`.
- Use `count=1` to disable test caching (`make test` does this).

## Domain Conventions

- **Value objects** validate on construction via `NewXxx()` — return `error` for invalid state. No setters.
- **Entity behavior** mutates state and emits events via `PullEvents()`. `ReconstituteTask()` / `ReconstituteDepartment()` are for DB reads only — does NOT emit events.
- **Status transitions**: `todo → in_progress → done` (reopenable: `done → in_progress`). Invalid transitions return `InvalidTransition` error.
- **`FindByID` returning nil** means not-found. Always check: `if t == nil { return nil, &NotFound{...} }`.

## Composition Root

`cmd/server/main.go` wires everything via `uber-go/fx`:
```go
fx.New(
    fx.Provide(
        loadConfig,
        newPool,
        postgres.NewTaskRepository,
        postgres.NewDepartmentRepository,
        taskapp.NewCreateTaskUseCase,  // repeat per use case
        taskhttp.NewTaskHandler,
        taskhttp.NewDepartmentHandler,
        taskhttp.NewAPIHandler,  // composes TaskHandler + DepartmentHandler
    ),
    fx.Invoke(startServer),
).Run()
```

Config from env vars (`PORT`, `DATABASE_URL`). No config file. Hardcoded defaults for local dev.

## Go Tooling

- **Requires Go 1.22+** (uses enhanced `http.ServeMux` with `{id}` path params and method routing). Go 1.26.3 in go.mod.
- **DI: `uber-go/fx`** — all wiring in `cmd/server/main.go`. No manual DI.
- **Query builder: `Masterminds/squirrel`** — PostgreSQL queries.
- Lint: `staticcheck ./...`
- Only test deps beyond std: `pgx/v5`, `uuid`, `oapi-codegen/runtime`, `testify`.

## Adding a Feature

1. Add/change OpenAPI spec → `make generate/openapi`
2. Add domain logic (entity behavior, value objects) → domain package
3. Add/update repository interface → domain package → `make generate/mocks`
4. Add use case → application package
5. Add/update handler → infrastructure/http → implements `ServerInterface`
6. Register constructor in `fx.Provide` in `cmd/server/main.go`
7. If new handler, wire into `APIHandler` composition in `cmd/server/main.go` and `api_handler.go`
8. Tests: domain (unit), application (mock), HTTP (httptest), integration (tagged)
