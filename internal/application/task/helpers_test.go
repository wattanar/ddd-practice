package taskapp_test

import (
	"context"
	"time"

	"github.com/wattanar/taskmanager/internal/domain/task"
	"github.com/wattanar/taskmanager/internal/domain/task/mocks"
)

var (
	ctx = context.Background()
	now = time.Now().UTC()
)

var _ task.TaskRepository = (*mocks.MockTaskRepository)(nil)
