package taskapp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wattanar/taskmanager/internal/application/task"
	"github.com/wattanar/taskmanager/internal/domain/task"
	"github.com/wattanar/taskmanager/internal/domain/task/mocks"
)

func TestListTasksUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewMockTaskRepository(t)

	status := task.TaskStatusTodo
	filter := task.TaskFilter{Status: &status}

	mockRepo.EXPECT().FindAll(ctx, filter).Return([]*task.Task{}, nil).Once()

	uc := taskapp.NewListTasksUseCase(mockRepo)

	result, err := uc.Execute(context.Background(), taskapp.ListTasksQuery{
		Status: &status,
	})

	require.NoError(t, err)
	assert.Empty(t, result)
}
