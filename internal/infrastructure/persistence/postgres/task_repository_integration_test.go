//go:build integration

package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wattanar/taskmanager/internal/domain/task"
	"github.com/wattanar/taskmanager/internal/infrastructure/persistence/postgres"
)

func TestTaskRepository_Integration(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	repo := postgres.NewTaskRepository(pool)

	t.Run("save and find by id", func(t *testing.T) {
		title, _ := task.NewTitle("Integration test")
		desc, _ := task.NewDescription("Testing the DB adapter")
		t1 := task.NewTask(title, desc, task.PriorityHigh)

		err := repo.Save(ctx, t1)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, t1.ID())
		require.NoError(t, err)
		require.NotNil(t, found)

		assert.Equal(t, t1.ID(), found.ID())
		assert.Equal(t, t1.Title().String(), found.Title().String())
		assert.Equal(t, t1.Status(), found.Status())
		assert.Equal(t, t1.Priority(), found.Priority())
	})

	t.Run("update existing task", func(t *testing.T) {
		title, _ := task.NewTitle("Update test")
		desc, _ := task.NewDescription("")
		t1 := task.NewTask(title, desc, task.PriorityLow)

		err := repo.Save(ctx, t1)
		require.NoError(t, err)

		newTitle, _ := task.NewTitle("Updated title")
		t1.UpdateTitle(newTitle)
		t1.UpdatePriority(task.PriorityHigh)

		err = repo.Save(ctx, t1)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, t1.ID())
		require.NoError(t, err)

		assert.Equal(t, "Updated title", found.Title().String())
		assert.Equal(t, task.PriorityHigh, found.Priority())
	})

	t.Run("find all with filter", func(t *testing.T) {
		t1 := task.NewTask(mustTitle("Filter A"), mustDesc(""), task.PriorityHigh)
		t1.ChangeStatus(task.TaskStatusInProgress)
		repo.Save(ctx, t1)

		t2 := task.NewTask(mustTitle("Filter B"), mustDesc(""), task.PriorityLow)
		repo.Save(ctx, t2)

		status := task.TaskStatusInProgress
		tasks, err := repo.FindAll(ctx, task.TaskFilter{Status: &status})
		require.NoError(t, err)

		for _, t := range tasks {
			assert.Equal(t, task.TaskStatusInProgress, t.Status())
		}
	})

	t.Run("delete task", func(t *testing.T) {
		t1 := task.NewTask(mustTitle("Delete me"), mustDesc(""), task.PriorityMedium)
		err := repo.Save(ctx, t1)
		require.NoError(t, err)

		err = repo.Delete(ctx, t1.ID())
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, t1.ID())
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("find by id not found", func(t *testing.T) {
		id := task.NewTaskID()
		found, err := repo.FindByID(ctx, id)
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func mustTitle(s string) task.Title {
	t, _ := task.NewTitle(s)
	return t
}

func mustDesc(s string) task.Description {
	d, _ := task.NewDescription(s)
	return d
}
