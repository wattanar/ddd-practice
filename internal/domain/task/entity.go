package task

import (
	"time"

	domainErrors "github.com/wattanar/taskmanager/internal/domain/errors"
)

type Task struct {
	id          TaskID
	title       Title
	description Description
	status      TaskStatus
	priority    Priority
	createdAt   time.Time
	updatedAt   time.Time
	events      []DomainEvent
}

func NewTask(title Title, description Description, priority Priority) *Task {
	now := time.Now()
	id := NewTaskID()
	t := &Task{
		id:        id,
		title:     title,
		status:    TaskStatusTodo,
		priority:  priority,
		createdAt: now,
		updatedAt: now,
	}
	if description.String() != "" {
		t.description = description
	}
	t.emit(NewTaskCreated(id, title.String()))
	return t
}

func ReconstituteTask(
	id TaskID,
	title Title,
	description Description,
	status TaskStatus,
	priority Priority,
	createdAt time.Time,
	updatedAt time.Time,
) *Task {
	return &Task{
		id:          id,
		title:       title,
		description: description,
		status:      status,
		priority:    priority,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (t *Task) ID() TaskID              { return t.id }
func (t *Task) Title() Title             { return t.title }
func (t *Task) Description() Description { return t.description }
func (t *Task) Status() TaskStatus       { return t.status }
func (t *Task) Priority() Priority       { return t.priority }
func (t *Task) CreatedAt() time.Time     { return t.createdAt }
func (t *Task) UpdatedAt() time.Time     { return t.updatedAt }

func (t *Task) UpdateTitle(title Title) {
	old := t.title.String()
	t.title = title
	t.updatedAt = time.Now()
	t.emit(NewTaskTitleChanged(t.id, old, title.String()))
}

func (t *Task) UpdateDescription(desc Description) {
	old := t.description.String()
	t.description = desc
	t.updatedAt = time.Now()
	t.emit(TaskDescriptionChanged{
		baseEvent:      baseEvent{occurredAt: time.Now()},
		TaskID:         t.id,
		OldDescription: old,
		NewDescription: desc.String(),
	})
}

func (t *Task) ChangeStatus(newStatus TaskStatus) error {
	if !t.status.CanTransitionTo(newStatus) {
		return &domainErrors.InvalidTransition{Current: t.status, Target: newStatus}
	}
	old := t.status
	t.status = newStatus
	t.updatedAt = time.Now()
	t.emit(NewTaskStatusChanged(t.id, old, newStatus))
	return nil
}

func (t *Task) UpdatePriority(priority Priority) {
	old := t.priority
	t.priority = priority
	t.updatedAt = time.Now()
	t.emit(TaskPriorityChanged{
		baseEvent:   baseEvent{occurredAt: time.Now()},
		TaskID:      t.id,
		OldPriority: old,
		NewPriority: priority,
	})
}

func (t *Task) Delete() {
	t.emit(NewTaskDeleted(t.id))
}

func (t *Task) emit(event DomainEvent) {
	t.events = append(t.events, event)
}

func (t *Task) PullEvents() []DomainEvent {
	events := t.events
	t.events = nil
	return events
}
