package deptapp

import (
	"context"
	"fmt"

	"github.com/wattanar/taskmanager/internal/domain/department"
	"github.com/wattanar/taskmanager/internal/domain/errors"
)

type DeleteDepartmentCommand struct {
	ID department.DepartmentID
}

type DeleteDepartmentUseCase struct {
	repo department.DepartmentRepository
}

func NewDeleteDepartmentUseCase(repo department.DepartmentRepository) *DeleteDepartmentUseCase {
	return &DeleteDepartmentUseCase{repo: repo}
}

func (uc *DeleteDepartmentUseCase) Execute(ctx context.Context, cmd DeleteDepartmentCommand) error {
	existing, err := uc.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("delete department: %w", err)
	}
	if existing == nil {
		return &errors.NotFound{Aggregate: "Department", ID: cmd.ID}
	}

	if err := uc.repo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("delete department: %w", err)
	}

	return nil
}
