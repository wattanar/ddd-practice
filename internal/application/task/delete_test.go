package taskapp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wattanar/taskmanager/internal/application/task"
	"github.com/wattanar/taskmanager/internal/domain/task"
	"github.com/wattanar/taskmanager/internal/domain/task/mocks"
)

func TestDeleteTaskUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewMockTaskRepository(t)

	id := task.NewTaskID()
	title, _ := task.NewTitle("Test")
	desc, _ := task.NewDescription("")
	existing := task.ReconstituteTask(id, title, desc, task.TaskStatusTodo, task.PriorityMedium, now, now)

	mockRepo.EXPECT().FindByID(ctx, id).Return(existing, nil).Once()
	mockRepo.EXPECT().Delete(ctx, id).Return(nil).Once()

	uc := taskapp.NewDeleteTaskUseCase(mockRepo)

	err := uc.Execute(context.Background(), taskapp.DeleteTaskCommand{ID: id})
	require.NoError(t, err)
}

func TestDeleteTaskUseCase_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockTaskRepository(t)

	id := task.NewTaskID()
	mockRepo.EXPECT().FindByID(ctx, id).Return(nil, nil).Once()

	uc := taskapp.NewDeleteTaskUseCase(mockRepo)

	err := uc.Execute(context.Background(), taskapp.DeleteTaskCommand{ID: id})
	require.Error(t, err)
}
