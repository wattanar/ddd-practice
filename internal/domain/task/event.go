package task

import "time"

type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

type baseEvent struct {
	occurredAt time.Time
}

func (e baseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

type TaskCreated struct {
	baseEvent
	TaskID TaskID
	Title  string
}

func NewTaskCreated(id TaskID, title string) TaskCreated {
	return TaskCreated{
		baseEvent: baseEvent{occurredAt: time.Now()},
		TaskID:    id,
		Title:     title,
	}
}

func (e TaskCreated) EventName() string { return "task.created" }

type TaskTitleChanged struct {
	baseEvent
	TaskID  TaskID
	OldTitle string
	NewTitle string
}

func NewTaskTitleChanged(id TaskID, oldTitle, newTitle string) TaskTitleChanged {
	return TaskTitleChanged{
		baseEvent: baseEvent{occurredAt: time.Now()},
		TaskID:    id,
		OldTitle:  oldTitle,
		NewTitle:  newTitle,
	}
}

func (e TaskTitleChanged) EventName() string { return "task.title_changed" }

type TaskDescriptionChanged struct {
	baseEvent
	TaskID        TaskID
	OldDescription string
	NewDescription string
}

func (e TaskDescriptionChanged) EventName() string { return "task.description_changed" }

type TaskStatusChanged struct {
	baseEvent
	TaskID    TaskID
	OldStatus TaskStatus
	NewStatus TaskStatus
}

func NewTaskStatusChanged(id TaskID, oldStatus, newStatus TaskStatus) TaskStatusChanged {
	return TaskStatusChanged{
		baseEvent: baseEvent{occurredAt: time.Now()},
		TaskID:    id,
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}
}

func (e TaskStatusChanged) EventName() string { return "task.status_changed" }

type TaskPriorityChanged struct {
	baseEvent
	TaskID      TaskID
	OldPriority Priority
	NewPriority Priority
}

func (e TaskPriorityChanged) EventName() string { return "task.priority_changed" }

type TaskDeleted struct {
	baseEvent
	TaskID TaskID
}

func NewTaskDeleted(id TaskID) TaskDeleted {
	return TaskDeleted{
		baseEvent: baseEvent{occurredAt: time.Now()},
		TaskID:    id,
	}
}

func (e TaskDeleted) EventName() string { return "task.deleted" }
