package taskapp

import (
	"context"
	"fmt"

	"github.com/wattanar/taskmanager/internal/domain/errors"
	"github.com/wattanar/taskmanager/internal/domain/task"
)

type UpdateTaskCommand struct {
	ID          task.TaskID
	Title       *string
	Description *string
	Status      *task.TaskStatus
	Priority    *task.Priority
}

type UpdateTaskUseCase struct {
	repo task.TaskRepository
}

func NewUpdateTaskUseCase(repo task.TaskRepository) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{repo: repo}
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, cmd UpdateTaskCommand) (*task.Task, error) {
	t, err := uc.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	if t == nil {
		return nil, &errors.NotFound{Aggregate: "Task", ID: cmd.ID}
	}

	if cmd.Title != nil {
		title, err := task.NewTitle(*cmd.Title)
		if err != nil {
			return nil, err
		}
		t.UpdateTitle(title)
	}

	if cmd.Description != nil {
		desc, err := task.NewDescription(*cmd.Description)
		if err != nil {
			return nil, err
		}
		t.UpdateDescription(desc)
	}

	if cmd.Status != nil {
		if err := t.ChangeStatus(*cmd.Status); err != nil {
			return nil, err
		}
	}

	if cmd.Priority != nil {
		t.UpdatePriority(*cmd.Priority)
	}

	if err := uc.repo.Save(ctx, t); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	return t, nil
}
