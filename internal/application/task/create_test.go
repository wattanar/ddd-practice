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

func TestCreateTaskUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewMockTaskRepository(t)
	mockRepo.EXPECT().Save(ctx, mock.AnythingOfType("*task.Task")).Return(nil).Once()

	uc := taskapp.NewCreateTaskUseCase(mockRepo)

	result, err := uc.Execute(context.Background(), taskapp.CreateTaskCommand{
		Title:       "Buy groceries",
		Description: "Milk and eggs",
		Priority:    task.PriorityHigh,
	})

	require.NoError(t, err)
	assert.Equal(t, "Buy groceries", result.Title().String())
	assert.Equal(t, "Milk and eggs", result.Description().String())
	assert.Equal(t, task.TaskStatusTodo, result.Status())
	assert.Equal(t, task.PriorityHigh, result.Priority())
}

func TestCreateTaskUseCase_EmptyTitle(t *testing.T) {
	mockRepo := mocks.NewMockTaskRepository(t)

	uc := taskapp.NewCreateTaskUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), taskapp.CreateTaskCommand{
		Title:       "",
		Description: "",
		Priority:    task.PriorityMedium,
	})

	require.Error(t, err)
}
