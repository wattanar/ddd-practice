package deptapp

import (
	"context"
	"fmt"

	"github.com/wattanar/taskmanager/internal/domain/department"
)

type CreateDepartmentCommand struct {
	Name        string
	Description string
}

type CreateDepartmentUseCase struct {
	repo department.DepartmentRepository
}

func NewCreateDepartmentUseCase(repo department.DepartmentRepository) *CreateDepartmentUseCase {
	return &CreateDepartmentUseCase{repo: repo}
}

func (uc *CreateDepartmentUseCase) Execute(ctx context.Context, cmd CreateDepartmentCommand) (*department.Department, error) {
	name, err := department.NewDepartmentName(cmd.Name)
	if err != nil {
		return nil, err
	}

	desc, err := department.NewDepartmentDescription(cmd.Description)
	if err != nil {
		return nil, err
	}

	d := department.NewDepartment(name, desc)

	if err := uc.repo.Save(ctx, d); err != nil {
		return nil, fmt.Errorf("create department: %w", err)
	}

	return d, nil
}
