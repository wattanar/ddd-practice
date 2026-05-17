package deptapp

import (
	"context"
	"fmt"

	"github.com/wattanar/taskmanager/internal/domain/department"
	"github.com/wattanar/taskmanager/internal/domain/errors"
)

type UpdateDepartmentCommand struct {
	ID          department.DepartmentID
	Name        *string
	Description *string
}

type UpdateDepartmentUseCase struct {
	repo department.DepartmentRepository
}

func NewUpdateDepartmentUseCase(repo department.DepartmentRepository) *UpdateDepartmentUseCase {
	return &UpdateDepartmentUseCase{repo: repo}
}

func (uc *UpdateDepartmentUseCase) Execute(ctx context.Context, cmd UpdateDepartmentCommand) (*department.Department, error) {
	d, err := uc.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("update department: %w", err)
	}
	if d == nil {
		return nil, &errors.NotFound{Aggregate: "Department", ID: cmd.ID}
	}

	if cmd.Name != nil {
		name, err := department.NewDepartmentName(*cmd.Name)
		if err != nil {
			return nil, err
		}
		d.UpdateName(name)
	}

	if cmd.Description != nil {
		desc, err := department.NewDepartmentDescription(*cmd.Description)
		if err != nil {
			return nil, err
		}
		d.UpdateDescription(desc)
	}

	if err := uc.repo.Save(ctx, d); err != nil {
		return nil, fmt.Errorf("update department: %w", err)
	}

	return d, nil
}
