package taskapp

import (
	"context"
	"fmt"

	"github.com/wattanar/taskmanager/internal/domain/errors"
	"github.com/wattanar/taskmanager/internal/domain/task"
)

type GetTaskQuery struct {
	ID task.TaskID
}

type GetTaskUseCase struct {
	repo task.TaskRepository
}

func NewGetTaskUseCase(repo task.TaskRepository) *GetTaskUseCase {
	return &GetTaskUseCase{repo: repo}
}

func (uc *GetTaskUseCase) Execute(ctx context.Context, q GetTaskQuery) (*task.Task, error) {
	t, err := uc.repo.FindByID(ctx, q.ID)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	if t == nil {
		return nil, &errors.NotFound{Aggregate: "Task", ID: q.ID}
	}
	return t, nil
}
