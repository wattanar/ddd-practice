package taskapp

import (
	"context"
	"fmt"

	"github.com/wattanar/taskmanager/internal/domain/task"
)

type ListTasksQuery struct {
	Status   *task.TaskStatus
	Priority *task.Priority
}

type ListTasksUseCase struct {
	repo task.TaskRepository
}

func NewListTasksUseCase(repo task.TaskRepository) *ListTasksUseCase {
	return &ListTasksUseCase{repo: repo}
}

func (uc *ListTasksUseCase) Execute(ctx context.Context, q ListTasksQuery) ([]*task.Task, error) {
	filter := task.TaskFilter{
		Status:   q.Status,
		Priority: q.Priority,
	}

	tasks, err := uc.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	return tasks, nil
}
