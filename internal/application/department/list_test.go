package deptapp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wattanar/taskmanager/internal/application/department"
	"github.com/wattanar/taskmanager/internal/domain/department"
	"github.com/wattanar/taskmanager/internal/domain/department/mocks"
)

func TestListDepartmentsUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)
	mockRepo.EXPECT().FindAll(ctx).Return([]*department.Department{}, nil).Once()

	uc := deptapp.NewListDepartmentsUseCase(mockRepo)

	result, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.Empty(t, result)
}
