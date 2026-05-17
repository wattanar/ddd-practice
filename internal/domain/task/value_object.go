package task

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/wattanar/taskmanager/internal/domain/errors"
)

type TaskID struct {
	uuid.UUID
}

func NewTaskID() TaskID {
	return TaskID{UUID: uuid.New()}
}

func ParseTaskID(s string) (TaskID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return TaskID{}, fmt.Errorf("parse task id: %w", err)
	}
	return TaskID{UUID: id}, nil
}

func (id TaskID) String() string { return id.UUID.String() }

type Title struct {
	value string
}

func NewTitle(s string) (Title, error) {
	if len(s) == 0 {
		return Title{}, &errors.InvalidArgument{Field: "title", Reason: "must not be empty"}
	}
	if len(s) > 200 {
		return Title{}, &errors.InvalidArgument{Field: "title", Reason: "must not exceed 200 characters"}
	}
	return Title{value: s}, nil
}

func (t Title) String() string {
	return t.value
}

type Description struct {
	value string
}

func NewDescription(s string) (Description, error) {
	if len(s) > 2000 {
		return Description{}, &errors.InvalidArgument{Field: "description", Reason: "must not exceed 2000 characters"}
	}
	return Description{value: s}, nil
}

func (d Description) String() string {
	return d.value
}

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

var validTransitions = map[TaskStatus][]TaskStatus{
	TaskStatusTodo:       {TaskStatusInProgress},
	TaskStatusInProgress: {TaskStatusDone},
	TaskStatusDone:       {TaskStatusInProgress},
}

func (s TaskStatus) CanTransitionTo(target TaskStatus) bool {
	for _, t := range validTransitions[s] {
		if t == target {
			return true
		}
	}
	return false
}

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)


