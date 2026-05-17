package deptapp_test

import (
	"context"
	"time"

	"github.com/wattanar/taskmanager/internal/domain/department"
	"github.com/wattanar/taskmanager/internal/domain/department/mocks"
)

var (
	ctx = context.Background()
	now = time.Now().UTC()
)

var _ department.DepartmentRepository = (*mocks.MockDepartmentRepository)(nil)
