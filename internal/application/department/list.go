package deptapp

import (
	"context"
	"fmt"

	"github.com/wattanar/taskmanager/internal/domain/department"
)

type ListDepartmentsUseCase struct {
	repo department.DepartmentRepository
}

func NewListDepartmentsUseCase(repo department.DepartmentRepository) *ListDepartmentsUseCase {
	return &ListDepartmentsUseCase{repo: repo}
}

func (uc *ListDepartmentsUseCase) Execute(ctx context.Context) ([]*department.Department, error) {
	departments, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list departments: %w", err)
	}
	return departments, nil
}
