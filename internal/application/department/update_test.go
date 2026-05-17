package deptapp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wattanar/taskmanager/internal/application/department"
	"github.com/wattanar/taskmanager/internal/domain/department"
	"github.com/wattanar/taskmanager/internal/domain/department/mocks"
)

func TestUpdateDepartmentUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)

	id := department.NewDepartmentID()
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("Original")
	existing := department.ReconstituteDepartment(id, name, desc, now, now)

	newName := "Engineering II"
	newDesc := "Updated description"

	mockRepo.EXPECT().FindByID(ctx, id).Return(existing, nil).Once()
	mockRepo.EXPECT().Save(ctx, mock.AnythingOfType("*department.Department")).
		Return(nil).
		Run(func(ctx context.Context, d *department.Department) {
			assert.Equal(t, newName, d.Name().String())
			assert.Equal(t, newDesc, d.Description().String())
		}).
		Once()

	uc := deptapp.NewUpdateDepartmentUseCase(mockRepo)

	result, err := uc.Execute(context.Background(), deptapp.UpdateDepartmentCommand{
		ID:          id,
		Name:        &newName,
		Description: &newDesc,
	})

	require.NoError(t, err)
	assert.Equal(t, newName, result.Name().String())
	assert.Equal(t, newDesc, result.Description().String())
}

func TestUpdateDepartmentUseCase_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)

	id := department.NewDepartmentID()
	mockRepo.EXPECT().FindByID(ctx, id).Return(nil, nil).Once()

	uc := deptapp.NewUpdateDepartmentUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), deptapp.UpdateDepartmentCommand{
		ID: id,
	})

	require.Error(t, err)
}
