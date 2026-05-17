package task_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainErrors "github.com/wattanar/taskmanager/internal/domain/errors"
	"github.com/wattanar/taskmanager/internal/domain/task"
)

func TestNewTitle(t *testing.T) {
	t.Run("valid title", func(t *testing.T) {
		title, err := task.NewTitle("Buy groceries")
		require.NoError(t, err)
		assert.Equal(t, "Buy groceries", title.String())
	})

	t.Run("empty title", func(t *testing.T) {
		_, err := task.NewTitle("")
		require.Error(t, err)
		var invalidArg *domainErrors.InvalidArgument
		assert.ErrorAs(t, err, &invalidArg)
	})

	t.Run("title too long", func(t *testing.T) {
		long := string(make([]byte, 201))
		_, err := task.NewTitle(long)
		require.Error(t, err)
		var invalidArg *domainErrors.InvalidArgument
		assert.ErrorAs(t, err, &invalidArg)
	})
}

func TestNewDescription(t *testing.T) {
	t.Run("valid description", func(t *testing.T) {
		desc, err := task.NewDescription("Buy milk and eggs")
		require.NoError(t, err)
		assert.Equal(t, "Buy milk and eggs", desc.String())
	})

	t.Run("empty description", func(t *testing.T) {
		desc, err := task.NewDescription("")
		require.NoError(t, err)
		assert.Equal(t, "", desc.String())
	})

	t.Run("description too long", func(t *testing.T) {
		long := string(make([]byte, 2001))
		_, err := task.NewDescription(long)
		require.Error(t, err)
	})
}

func TestTaskStatusTransitions(t *testing.T) {
	tests := []struct {
		name     string
		current  task.TaskStatus
		target   task.TaskStatus
		expected bool
	}{
		{"todo → in_progress", task.TaskStatusTodo, task.TaskStatusInProgress, true},
		{"todo → done", task.TaskStatusTodo, task.TaskStatusDone, false},
		{"in_progress → done", task.TaskStatusInProgress, task.TaskStatusDone, true},
		{"in_progress → todo", task.TaskStatusInProgress, task.TaskStatusTodo, false},
		{"done → in_progress", task.TaskStatusDone, task.TaskStatusInProgress, true},
		{"done → todo", task.TaskStatusDone, task.TaskStatusTodo, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.current.CanTransitionTo(tt.target))
		})
	}
}

func TestParseTaskID(t *testing.T) {
	t.Run("valid uuid", func(t *testing.T) {
		id, err := task.ParseTaskID("550e8400-e29b-41d4-a716-446655440000")
		require.NoError(t, err)
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
	})

	t.Run("invalid uuid", func(t *testing.T) {
		_, err := task.ParseTaskID("not-a-uuid")
		require.Error(t, err)
	})
}
