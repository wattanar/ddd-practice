package task_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wattanar/taskmanager/internal/domain/task"
)

func TestNewTask(t *testing.T) {
	title, _ := task.NewTitle("Test task")
	desc, _ := task.NewDescription("A description")
	t1 := task.NewTask(title, desc, task.PriorityHigh)

	assert.Equal(t, "Test task", t1.Title().String())
	assert.Equal(t, "A description", t1.Description().String())
	assert.Equal(t, task.TaskStatusTodo, t1.Status())
	assert.Equal(t, task.PriorityHigh, t1.Priority())
	assert.False(t, t1.CreatedAt().IsZero())
	assert.False(t, t1.UpdatedAt().IsZero())

	events := t1.PullEvents()
	assert.Len(t, events, 1)
	assert.Equal(t, "task.created", events[0].EventName())
}

func TestTask_ChangeStatus(t *testing.T) {
	title, _ := task.NewTitle("Test")
	desc, _ := task.NewDescription("")
	t1 := task.NewTask(title, desc, task.PriorityMedium)
	t1.PullEvents()

	t.Run("valid transition", func(t *testing.T) {
		err := t1.ChangeStatus(task.TaskStatusInProgress)
		require.NoError(t, err)
		assert.Equal(t, task.TaskStatusInProgress, t1.Status())

		events := t1.PullEvents()
		assert.Len(t, events, 1)
		statusEvent, ok := events[0].(task.TaskStatusChanged)
		assert.True(t, ok)
		assert.Equal(t, task.TaskStatusTodo, statusEvent.OldStatus)
		assert.Equal(t, task.TaskStatusInProgress, statusEvent.NewStatus)
	})

	t.Run("invalid transition", func(t *testing.T) {
		err := t1.ChangeStatus(task.TaskStatusTodo)
		require.Error(t, err)
		assert.Equal(t, task.TaskStatusInProgress, t1.Status())
	})
}

func TestTask_UpdateTitle(t *testing.T) {
	title, _ := task.NewTitle("Original")
	desc, _ := task.NewDescription("")
	t1 := task.NewTask(title, desc, task.PriorityMedium)
	t1.PullEvents()

	newTitle, _ := task.NewTitle("Updated")
	t1.UpdateTitle(newTitle)

	assert.Equal(t, "Updated", t1.Title().String())

	events := t1.PullEvents()
	assert.Len(t, events, 1)
	assert.Equal(t, "task.title_changed", events[0].EventName())
}

func TestTask_UpdateDescription(t *testing.T) {
	title, _ := task.NewTitle("Test")
	desc, _ := task.NewDescription("Original desc")
	t1 := task.NewTask(title, desc, task.PriorityLow)
	t1.PullEvents()

	newDesc, _ := task.NewDescription("Updated desc")
	t1.UpdateDescription(newDesc)

	assert.Equal(t, "Updated desc", t1.Description().String())
}

func TestTask_UpdatePriority(t *testing.T) {
	title, _ := task.NewTitle("Test")
	desc, _ := task.NewDescription("")
	t1 := task.NewTask(title, desc, task.PriorityLow)
	t1.PullEvents()

	t1.UpdatePriority(task.PriorityCritical)
	assert.Equal(t, task.PriorityCritical, t1.Priority())
}

func TestTask_PullEvents(t *testing.T) {
	title, _ := task.NewTitle("Test")
	desc, _ := task.NewDescription("")
	t1 := task.NewTask(title, desc, task.PriorityMedium)

	events := t1.PullEvents()
	assert.Len(t, events, 1)

	empty := t1.PullEvents()
	assert.Len(t, empty, 0)
}

func TestTask_Timestamps(t *testing.T) {
	title, _ := task.NewTitle("Test")
	desc, _ := task.NewDescription("")
	t1 := task.NewTask(title, desc, task.PriorityMedium)

	originalUpdated := t1.UpdatedAt()

	time.Sleep(time.Millisecond)

	newTitle, _ := task.NewTitle("Updated")
	t1.UpdateTitle(newTitle)

	assert.True(t, t1.UpdatedAt().After(originalUpdated))
}

func TestReconstituteTask(t *testing.T) {
	now := time.Now()
	id := task.NewTaskID()
	title, _ := task.NewTitle("Reconstituted")
	desc, _ := task.NewDescription("Desc")

	t1 := task.ReconstituteTask(id, title, desc, task.TaskStatusInProgress, task.PriorityHigh, now, now)

	assert.Equal(t, id, t1.ID())
	assert.Equal(t, "Reconstituted", t1.Title().String())
	assert.Equal(t, task.TaskStatusInProgress, t1.Status())
	assert.Equal(t, task.PriorityHigh, t1.Priority())
	assert.Equal(t, now, t1.CreatedAt())
	assert.Equal(t, now, t1.UpdatedAt())

	events := t1.PullEvents()
	assert.Len(t, events, 0)
}
