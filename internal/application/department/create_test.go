package deptapp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wattanar/taskmanager/internal/application/department"
	"github.com/wattanar/taskmanager/internal/domain/department/mocks"
)

func TestCreateDepartmentUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)
	mockRepo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*department.Department")).Return(nil).Once()

	uc := deptapp.NewCreateDepartmentUseCase(mockRepo)

	result, err := uc.Execute(context.Background(), deptapp.CreateDepartmentCommand{
		Name:        "Engineering",
		Description: "Engineering department",
	})

	require.NoError(t, err)
	assert.Equal(t, "Engineering", result.Name().String())
	assert.Equal(t, "Engineering department", result.Description().String())
}

func TestCreateDepartmentUseCase_EmptyName(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)

	uc := deptapp.NewCreateDepartmentUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), deptapp.CreateDepartmentCommand{
		Name: "",
	})

	require.Error(t, err)
}
