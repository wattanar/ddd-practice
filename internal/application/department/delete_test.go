package deptapp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wattanar/taskmanager/internal/application/department"
	"github.com/wattanar/taskmanager/internal/domain/department"
	"github.com/wattanar/taskmanager/internal/domain/department/mocks"
)

func TestDeleteDepartmentUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)

	id := department.NewDepartmentID()
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("")
	existing := department.ReconstituteDepartment(id, name, desc, now, now)

	mockRepo.EXPECT().FindByID(ctx, id).Return(existing, nil).Once()
	mockRepo.EXPECT().Delete(ctx, id).Return(nil).Once()

	uc := deptapp.NewDeleteDepartmentUseCase(mockRepo)

	err := uc.Execute(context.Background(), deptapp.DeleteDepartmentCommand{ID: id})
	require.NoError(t, err)
}

func TestDeleteDepartmentUseCase_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)

	id := department.NewDepartmentID()
	mockRepo.EXPECT().FindByID(ctx, id).Return(nil, nil).Once()

	uc := deptapp.NewDeleteDepartmentUseCase(mockRepo)

	err := uc.Execute(context.Background(), deptapp.DeleteDepartmentCommand{ID: id})
	require.Error(t, err)
}
