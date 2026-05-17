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

func TestGetDepartmentUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)

	id := department.NewDepartmentID()
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("Desc")
	expected := department.ReconstituteDepartment(id, name, desc, now, now)

	mockRepo.EXPECT().FindByID(ctx, id).Return(expected, nil).Once()

	uc := deptapp.NewGetDepartmentUseCase(mockRepo)

	result, err := uc.Execute(context.Background(), deptapp.GetDepartmentQuery{ID: id})

	require.NoError(t, err)
	assert.Equal(t, id, result.ID())
	assert.Equal(t, "Engineering", result.Name().String())
}

func TestGetDepartmentUseCase_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)

	id := department.NewDepartmentID()
	mockRepo.EXPECT().FindByID(ctx, id).Return(nil, nil).Once()

	uc := deptapp.NewGetDepartmentUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), deptapp.GetDepartmentQuery{ID: id})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
