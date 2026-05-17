# Tutorial — DDD in Practice

This tutorial walks through how Clean Architecture + DDD + Hexagonal patterns are applied in this project, how requests flow through the layers, and how to add new features.

## Table of Contents

1. [Request Lifecycle](#1-request-lifecycle)
2. [Layers Explained](#2-layers-explained)
3. [DDD Building Blocks in Code](#3-ddd-building-blocks-in-code)
4. [Adding a New Feature — Full Walkthrough](#4-adding-a-new-feature--full-walkthrough)
5. [Adding a New Bounded Context](#5-adding-a-new-bounded-context)
6. [How to Think in DDD](#6-how-to-think-in-ddd)

---

## 1. Request Lifecycle

Every API request follows this path:

```
HTTP Request
    │
    ▼
┌─────────────────────────────────────┐
│ infrastructure/http                 │
│   handler.go (ServerInterface impl) │  ← Translates OpenAPI DTO → domain types
│   - decodes JSON                    │     e.g., spec.Task → domain/task.Task
│   - calls use case                  │
│   - encodes domain → JSON response  │
└─────────────┬───────────────────────┘
              │ taskapp.CreateTaskCommand
              ▼
┌─────────────────────────────────────┐
│ application/task                    │
│   create.go (use case)              │  ← Orchestrates: validate → domain logic → save
│   - takes Command struct            │
│   - constructs domain objects       │
│   - calls repository.Save()         │
└─────────────┬───────────────────────┘
              │ domain.TaskRepository.Save(task)
              ▼
┌─────────────────────────────────────┐
│ domain/task                         │
│   entity.go                         │  ← Pure business rules
│   value_object.go                   │     (no imports from pgx, http, etc.)
│   event.go                          │
│   repository.go (PORT interface)    │
└─────────────┬───────────────────────┘
              │ interface implemented by:
              ▼
┌─────────────────────────────────────┐
│ infrastructure/persistence/postgres │
│   task_repository.go                │  ← Implements port using pgx + squirrel
└─────────────────────────────────────┘
```

**Key rule:** Each layer only knows about the layer directly inside it. The HTTP handler never touches the database. The domain never imports `pgx`, `squirrel`, or any framework.

---

## 2. Layers Explained

### Domain (`internal/domain/`)

**Zero external dependencies.** Pure Go standard library only. This is the heart of the application.

What lives here:

| File | Purpose |
|---|---|
| `task/value_object.go` | `TaskID`, `Title`, `Description`, `TaskStatus`, `Priority` — validated on construction |
| `task/entity.go` | `Task` aggregate — behavior methods + domain event emission |
| `task/event.go` | `TaskCreated`, `TaskStatusChanged`, etc. — record of what happened |
| `task/repository.go` | `TaskRepository` interface — the **port** that infrastructure implements |
| `errors/errors.go` | `NotFound`, `InvalidArgument`, `InvalidTransition` — standard domain errors |

**The golden rule:** If you find yourself importing `pgx`, `net/http`, or any third-party library in domain, you've violated the architecture. Stop and move that code to infrastructure.

### Application (`internal/application/`)

One file per use case. Depends only on `domain/task`. Each use case is a struct with a single `Execute` method:

```go
type CreateTaskUseCase struct {
    repo task.TaskRepository  // ← depends on interface, not implementation
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, cmd CreateTaskCommand) (*task.Task, error) {
    // 1. Validate input via value object constructors
    title, _ := task.NewTitle(cmd.Title)
    // 2. Call domain behavior
    t := task.NewTask(title, desc, priority)
    // 3. Persist via port
    uc.repo.Save(ctx, t)
    return t, nil
}
```

Use cases never:
- Import infrastructure packages
- Open database connections
- Handle HTTP requests/responses

### Infrastructure (`internal/infrastructure/`)

Two kinds of adapters:

**Driven adapters** (called BY the application):
- `persistence/postgres/task_repository.go` — implements `task.TaskRepository` using `pgx` + `squirrel`

**Driving adapters** (call the application):
- `http/handler.go` — implements generated `ServerInterface`, translates HTTP JSON ↔ domain types

---

## 3. DDD Building Blocks in Code

### Value Objects

Immutable. Validated on construction. No setters.

```go
type Title struct{ value string }

func NewTitle(s string) (Title, error) {
    if len(s) == 0 {
        return Title{}, &errors.InvalidArgument{Field: "title", Reason: "must not be empty"}
    }
    if len(s) > 200 {
        return Title{}, &errors.InvalidArgument{Field: "title", Reason: "must not exceed 200 characters"}
    }
    return Title{value: s}, nil
}

func (t Title) String() string { return t.value }
```

**How to use:** `title, err := task.NewTitle("Buy milk")` — if validation fails, you get a domain error.

### Entity (Aggregate Root)

Has identity (ID). Behavior methods mutate state and emit events.

```go
t := task.NewTask(title, desc, task.PriorityHigh)
// t.Status() == TaskStatusTodo

err := t.ChangeStatus(task.TaskStatusInProgress)
// t.Status() == TaskStatusInProgress

events := t.PullEvents()
// events[0].EventName() == "task.created"
// events[1].EventName() == "task.status_changed"
```

**Two creation paths:**

| Path | When | Events? |
|---|---|---|
| `NewTask()` | Creating a new task | Emits `TaskCreated` |
| `ReconstituteTask()` | Loading from DB | No events (read-only reconstruction) |

### Domain Events

Events are collected on the entity and pulled after the use case completes:

```go
// In use case:
t := task.NewTask(...)
repo.Save(ctx, t)
events := t.PullEvents()  // ← now publish to message queue, etc.
```

Currently events are collected but not published. The pattern is ready for adding an event publisher in infrastructure.

### Repository (Port)

Defined in domain as an interface:

```go
type TaskRepository interface {
    Save(ctx context.Context, task *Task) error
    FindByID(ctx context.Context, id TaskID) (*Task, error)
    FindAll(ctx context.Context, filter TaskFilter) ([]*Task, error)
    Delete(ctx context.Context, id TaskID) error
}
```

Implemented in infrastructure:

```go
type TaskRepository struct {
    pool    *pgxpool.Pool
    builder sq.StatementBuilderType
}
```

**Critical convention:** `FindByID` returns `(nil, nil)` when a record is not found — this is **not** an error. Every use case must check for nil:

```go
t, err := repo.FindByID(ctx, id)
if err != nil { return nil, err }
if t == nil { return nil, &domainErrors.NotFound{Aggregate: "Task", ID: id} }
```

### Domain Errors

Three error types used across layers:

| Error | HTTP Status | When |
|---|---|---|
| `&NotFound{}` | 404 | Task not found |
| `&InvalidArgument{}` | 400 | Title empty, description too long |
| `&InvalidTransition{}` | 400 | Status change not allowed (e.g., `todo → done`) |

The HTTP handler maps these to HTTP status codes via `writeDomainError()`.

---

## 4. Adding a New Feature — Full Walkthrough

### Scenario: Add "due date" to tasks

Business requirement: "Tasks can have an optional due date. When a task passes its due date, it should be marked as `overdue`."

#### Step 1: Update the OpenAPI spec

`api/openapi.yaml`:

```yaml
Task:
  properties:
    dueDate:
      type: string
      format: date-time
      description: Optional due date

CreateTaskRequest:
  properties:
    dueDate:
      type: string
      format: date-time

UpdateTaskRequest:
  properties:
    dueDate:
      type: string
      format: date-time
```

Run: `make generate/openapi`

#### Step 2: Add domain value object

`internal/domain/task/value_object.go`:

```go
type DueDate struct {
    value time.Time
}

func NewDueDate(t time.Time) DueDate {
    return DueDate{value: t}
}

func (d DueDate) IsZero() bool  { return d.value.IsZero() }
func (d DueDate) Time() time.Time { return d.value }
func (d DueDate) IsOverdue() bool {
    return !d.value.IsZero() && time.Now().After(d.value)
}
```

#### Step 3: Add entity behavior

`internal/domain/task/entity.go`:

```go
type Task struct {
    // ... existing fields
    dueDate DueDate
}

// In NewTask():
func NewTask(...) *Task {
    // ... existing
    t.emit(NewTaskCreated(id, title.String()))
    return t
}

func (t *Task) SetDueDate(due DueDate) {
    t.dueDate = due
    t.updatedAt = time.Now()
    t.emit(TaskDueDateChanged{baseEvent: baseEvent{occurredAt: time.Now()}, TaskID: t.id, DueDate: due})
}

func (t *Task) DueDate() DueDate { return t.dueDate }

func (t *Task) CheckOverdue() {
    if t.dueDate.IsOverdue() {
        t.status = TaskStatusOverdue  // ← add TaskStatusOverdue to TaskStatus
    }
}
```

#### Step 4: Update repository interface

`internal/domain/task/repository.go` — add `due_date` to the filter if needed, or leave as-is since Save uses the full entity.

Run: `make generate/mocks`

#### Step 5: Update PostgreSQL repository

`internal/infrastructure/persistence/postgres/task_repository.go`:

```go
// In Save():
Columns("id", "title", "description", "status", "priority", "due_date", "created_at", "updated_at")
Values(t.ID().UUID, ..., t.DueDate().Time(), ...)
Suffix("ON CONFLICT (id) DO UPDATE SET ..., due_date = EXCLUDED.due_date, ...")

// In FindByID and FindAll:
// Add due_date to SELECT and scan
```

#### Step 5b: Add migration

`migrations/v002_add_due_date.sql`:

```sql
--liquibase formatted sql

--changeset app:002
ALTER TABLE tasks ADD COLUMN due_date TIMESTAMPTZ;
--rollback ALTER TABLE tasks DROP COLUMN due_date;
```

Add to `changelog-root.xml`:

```xml
<include file="v002_add_due_date.sql" relativeToChangelogFile="true"/>
```

Run: `make migrate`

#### Step 6: Update use case

`internal/application/task/create.go`:

```go
type CreateTaskCommand struct {
    Title       string
    Description string
    Priority    task.Priority
    DueDate     time.Time
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, cmd CreateTaskCommand) (*task.Task, error) {
    // ... existing validation ...
    t := task.NewTask(title, desc, cmd.Priority)
    if !cmd.DueDate.IsZero() {
        t.SetDueDate(task.NewDueDate(cmd.DueDate))
    }
    // ...
}
```

#### Step 7: Update HTTP handler

`internal/infrastructure/http/handler.go` — pass `DueDate` from the request DTO to the use case command.

#### Step 8: Add tests

| Test | Location | What to test |
|---|---|---|
| Value object | `domain/task/value_object_test.go` | `NewDueDate` validation |
| Entity behavior | `domain/task/entity_test.go` | `SetDueDate`, `CheckOverdue` |
| Use case | `application/task/create_test.go` | New field in command |
| HTTP | `http/handler_test.go` | `dueDate` in JSON request/response |

Run: `make test`

#### Step 9: Wire in main

Nothing extra needed — no new use case was added, just fields updated.

---

## 5. Adding a New Bounded Context

### Scenario: Add "Projects" that contain tasks

This means adding a whole new domain concept. Follow the same layered pattern:

```
internal/
├── domain/
│   └── project/
│       ├── entity.go           # Project aggregate
│       ├── value_object.go     # ProjectName, ProjectID
│       ├── event.go            # ProjectCreated, etc.
│       └── repository.go       # ProjectRepository port
├── application/
│   └── project/
│       ├── create.go
│       ├── get.go
│       └── list.go
└── infrastructure/
    ├── persistence/postgres/
    │   └── project_repository.go
    └── http/
        ├── project_handler.go   # New handler implementing ServerInterface
        └── router.go            # Register new routes
```

**API spec** — add new paths to `api/openapi.yaml`:

```yaml
paths:
  /api/v1/projects:
    get:
      operationId: listProjects
    post:
      operationId: createProject
  /api/v1/projects/{id}:
    get:
      operationId: getProject
    patch:
      operationId: updateProject
    delete:
      operationId: deleteProject
```

Then follow the same 9-step process from Section 4 for each endpoint.

**Cross-context communication** (when Project needs Task data):
- Project use case depends on `task.TaskRepository` (the port, not the implementation)
- Wire both repositories to the use case in `main.go`
- Never call infrastructure directly from a use case

```go
type ListProjectTasksUseCase struct {
    projectRepo project.ProjectRepository
    taskRepo    task.TaskRepository  // ← depends on port from another context
}
```

---

## 6. How to Think in DDD

### Ask These Questions First

When a new requirement arrives, work through these questions in order:

1. **What business operation is this?** Not "what data do I store" but "what action does the user take?"
   - ❌ "Add a notes column to the task table"
   - ✓ "The user can add free-form notes to a task"

2. **Where does the business logic live?**
   - Does the entity need a new behavior method? (`task.AddNote(note)`)
   - Does a new value object need validation rules? (`NewNote(s) error`)
   - Are there new state transitions? (`active → archived`)

3. **What changes at each layer?**

```
OpenAPI spec  →  Domain logic  →  Repository interface
                                        ↓
                                  Generated mocks
                                        ↓
                                  Application use case
                                        ↓
                                  HTTP handler (DTO translation)
                                        ↓
                                  Repository implementation (SQL)
                                        ↓
                                  Database migration
                                        ↓
                                  Tests at every layer
```

### Common Mistakes to Catch

| Mistake | How to Spot | Fix |
|---|---|---|
| Domain imports external lib | `import "pgx"` in `domain/` | Move to infrastructure |
| Handler calls repo directly | `handler.repo.FindByID()` | Go through use case |
| Anemic entity (getters only) | Entity has no behavior methods | Move logic into entity |
| Repository per table | `TaskItemRepository`, `TaskNoteRepository` | One per aggregate |
| Use case returns DTO | `CreateTaskUseCase` returns `spec.Task` | Return domain type, convert in handler |

### Testing Heuristic

> If you can test your domain logic with `go test` and no database running, your architecture is correct.

```
Domain tests:  go test ./internal/domain/...     # instant, no infra
App tests:     go test ./internal/application/... # mock repo, fast
HTTP tests:    go test ./internal/infrastructure/http/... # httptest, no DB
Integration:   go test -tags=integration ...      # needs real Postgres
```

### Preparing for New Requirements

Before coding a new feature, prepare:

1. **Review the OpenAPI spec** — is this a new endpoint or a field change?
2. **Check the domain** — does the concept already exist? Does the entity need a new behavior?
3. **Review the dependency flow** — does this require a new port? A new use case?
4. **Write the tests first** — domain tests define the behavior, HTTP tests define the API contract
5. **Commit generated code** — run `make generate` and commit the output so CI doesn't regenerate it differently
6. **Never break the dependency rule** — every time you add an import, check which layer it belongs to
