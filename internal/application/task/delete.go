package taskapp

import (
	"context"
	"fmt"

	"github.com/wattanar/taskmanager/internal/domain/errors"
	"github.com/wattanar/taskmanager/internal/domain/task"
)

type DeleteTaskCommand struct {
	ID task.TaskID
}

type DeleteTaskUseCase struct {
	repo task.TaskRepository
}

func NewDeleteTaskUseCase(repo task.TaskRepository) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{repo: repo}
}

func (uc *DeleteTaskUseCase) Execute(ctx context.Context, cmd DeleteTaskCommand) error {
	existing, err := uc.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if existing == nil {
		return &errors.NotFound{Aggregate: "Task", ID: cmd.ID}
	}

	if err := uc.repo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}
