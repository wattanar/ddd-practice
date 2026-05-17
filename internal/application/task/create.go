package taskapp

import (
	"context"
	"fmt"

	"github.com/wattanar/taskmanager/internal/domain/task"
)

type CreateTaskCommand struct {
	Title       string
	Description string
	Priority    task.Priority
}

type CreateTaskUseCase struct {
	repo task.TaskRepository
}

func NewCreateTaskUseCase(repo task.TaskRepository) *CreateTaskUseCase {
	return &CreateTaskUseCase{repo: repo}
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, cmd CreateTaskCommand) (*task.Task, error) {
	title, err := task.NewTitle(cmd.Title)
	if err != nil {
		return nil, err
	}

	desc, err := task.NewDescription(cmd.Description)
	if err != nil {
		return nil, err
	}

	if cmd.Priority == "" {
		cmd.Priority = task.PriorityMedium
	}

	t := task.NewTask(title, desc, cmd.Priority)

	if err := uc.repo.Save(ctx, t); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return t, nil
}
