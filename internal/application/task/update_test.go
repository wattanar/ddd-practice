package taskapp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wattanar/taskmanager/internal/application/task"
	"github.com/wattanar/taskmanager/internal/domain/task"
	"github.com/wattanar/taskmanager/internal/domain/task/mocks"
)

func TestUpdateTaskUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewMockTaskRepository(t)

	id := task.NewTaskID()
	title, _ := task.NewTitle("Original title")
	desc, _ := task.NewDescription("Original desc")
	existing := task.ReconstituteTask(id, title, desc, task.TaskStatusTodo, task.PriorityMedium, now, now)

	newTitle := "Updated title"
	newStatus := task.TaskStatusInProgress

	mockRepo.EXPECT().FindByID(ctx, id).Return(existing, nil).Once()
	mockRepo.EXPECT().Save(ctx, mock.AnythingOfType("*task.Task")).
		Return(nil).
		Run(func(ctx context.Context, saved *task.Task) {
			assert.Equal(t, newTitle, saved.Title().String())
			assert.Equal(t, task.TaskStatusInProgress, saved.Status())
		}).
		Once()

	uc := taskapp.NewUpdateTaskUseCase(mockRepo)

	result, err := uc.Execute(context.Background(), taskapp.UpdateTaskCommand{
		ID:     id,
		Title:  &newTitle,
		Status: &newStatus,
	})

	require.NoError(t, err)
	assert.Equal(t, newTitle, result.Title().String())
	assert.Equal(t, task.TaskStatusInProgress, result.Status())
}

func TestUpdateTaskUseCase_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockTaskRepository(t)

	id := task.NewTaskID()
	mockRepo.EXPECT().FindByID(ctx, id).Return(nil, nil).Once()

	uc := taskapp.NewUpdateTaskUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), taskapp.UpdateTaskCommand{
		ID: id,
	})

	require.Error(t, err)
}
