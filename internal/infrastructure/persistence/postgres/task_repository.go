package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wattanar/taskmanager/internal/domain/task"
)

type TaskRepository struct {
	pool   *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *TaskRepository) Save(ctx context.Context, t *task.Task) error {
	query, args, err := r.builder.Insert("tasks").
		Columns("id", "title", "description", "status", "priority", "created_at", "updated_at").
		Values(t.ID().UUID, t.Title().String(), t.Description().String(), string(t.Status()), string(t.Priority()), t.CreatedAt(), t.UpdatedAt()).
		Suffix("ON CONFLICT (id) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description, status = EXCLUDED.status, priority = EXCLUDED.priority, updated_at = EXCLUDED.updated_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("build save query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("save task: %w", err)
	}

	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id task.TaskID) (*task.Task, error) {
	query, args, err := r.builder.Select("id", "title", "description", "status", "priority", "created_at", "updated_at").
		From("tasks").
		Where(sq.Eq{"id": id.UUID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build find by id query: %w", err)
	}

	row := r.pool.QueryRow(ctx, query, args...)

	var (
		taskID      string
		title       string
		description string
		status      string
		priority    string
		createdAt   time.Time
		updatedAt   time.Time
	)

	err = row.Scan(&taskID, &title, &description, &status, &priority, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find task by id: %w", err)
	}

	return r.mapToTask(taskID, title, description, status, priority, createdAt, updatedAt)
}

func (r *TaskRepository) FindAll(ctx context.Context, filter task.TaskFilter) ([]*task.Task, error) {
	b := r.builder.Select("id", "title", "description", "status", "priority", "created_at", "updated_at").
		From("tasks").
		OrderBy("created_at DESC")

	if filter.Status != nil {
		b = b.Where(sq.Eq{"status": string(*filter.Status)})
	}

	if filter.Priority != nil {
		b = b.Where(sq.Eq{"priority": string(*filter.Priority)})
	}

	query, args, err := b.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build find all query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("find all tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		var (
			taskID      string
			title       string
			description string
			status      string
			priority    string
			createdAt   time.Time
			updatedAt   time.Time
		)

		if err := rows.Scan(&taskID, &title, &description, &status, &priority, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan task row: %w", err)
		}

		t, err := r.mapToTask(taskID, title, description, status, priority, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate task rows: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id task.TaskID) error {
	query, args, err := r.builder.Delete("tasks").
		Where(sq.Eq{"id": id.UUID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}

func (r *TaskRepository) mapToTask(
	taskID, title, description, status, priority string,
	createdAt, updatedAt time.Time,
) (*task.Task, error) {
	id, err := task.ParseTaskID(taskID)
	if err != nil {
		return nil, err
	}

	t, err := task.NewTitle(title)
	if err != nil {
		return nil, err
	}

	d, err := task.NewDescription(description)
	if err != nil {
		return nil, err
	}

	return task.ReconstituteTask(
		id,
		t,
		d,
		task.TaskStatus(status),
		task.Priority(priority),
		createdAt,
		updatedAt,
	), nil
}
