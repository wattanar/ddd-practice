package deptapp

import (
	"context"
	"fmt"

	"github.com/wattanar/taskmanager/internal/domain/department"
	"github.com/wattanar/taskmanager/internal/domain/errors"
)

type GetDepartmentQuery struct {
	ID department.DepartmentID
}

type GetDepartmentUseCase struct {
	repo department.DepartmentRepository
}

func NewGetDepartmentUseCase(repo department.DepartmentRepository) *GetDepartmentUseCase {
	return &GetDepartmentUseCase{repo: repo}
}

func (uc *GetDepartmentUseCase) Execute(ctx context.Context, q GetDepartmentQuery) (*department.Department, error) {
	d, err := uc.repo.FindByID(ctx, q.ID)
	if err != nil {
		return nil, fmt.Errorf("get department: %w", err)
	}
	if d == nil {
		return nil, &errors.NotFound{Aggregate: "Department", ID: q.ID}
	}
	return d, nil
}
