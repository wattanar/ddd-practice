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

func TestGetTaskUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewMockTaskRepository(t)

	id := task.NewTaskID()
	title, _ := task.NewTitle("Test task")
	desc, _ := task.NewDescription("Description")
	expected := task.ReconstituteTask(id, title, desc, task.TaskStatusTodo, task.PriorityMedium, now, now)

	mockRepo.EXPECT().FindByID(ctx, id).Return(expected, nil).Once()

	uc := taskapp.NewGetTaskUseCase(mockRepo)

	result, err := uc.Execute(context.Background(), taskapp.GetTaskQuery{ID: id})

	require.NoError(t, err)
	assert.Equal(t, id, result.ID())
	assert.Equal(t, "Test task", result.Title().String())
}

func TestGetTaskUseCase_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockTaskRepository(t)

	id := task.NewTaskID()
	mockRepo.EXPECT().FindByID(ctx, id).Return(nil, nil).Once()

	uc := taskapp.NewGetTaskUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), taskapp.GetTaskQuery{ID: id})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
