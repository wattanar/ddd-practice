package task

import "context"

type TaskFilter struct {
	Status   *TaskStatus
	Priority *Priority
}

type TaskRepository interface {
	Save(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id TaskID) (*Task, error)
	FindAll(ctx context.Context, filter TaskFilter) ([]*Task, error)
	Delete(ctx context.Context, id TaskID) error
}
