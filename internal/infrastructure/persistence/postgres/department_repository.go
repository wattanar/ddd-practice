package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wattanar/taskmanager/internal/domain/department"
)

type DepartmentRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewDepartmentRepository(pool *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *DepartmentRepository) Save(ctx context.Context, d *department.Department) error {
	query, args, err := r.builder.Insert("departments").
		Columns("id", "name", "description", "created_at", "updated_at").
		Values(d.ID().UUID, d.Name().String(), d.Description().String(), d.CreatedAt(), d.UpdatedAt()).
		Suffix("ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, updated_at = EXCLUDED.updated_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("build save department query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("save department: %w", err)
	}

	return nil
}

func (r *DepartmentRepository) FindByID(ctx context.Context, id department.DepartmentID) (*department.Department, error) {
	query, args, err := r.builder.Select("id", "name", "description", "created_at", "updated_at").
		From("departments").
		Where(sq.Eq{"id": id.UUID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build find department by id query: %w", err)
	}

	row := r.pool.QueryRow(ctx, query, args...)

	var (
		deptID      string
		name        string
		description string
		createdAt   time.Time
		updatedAt   time.Time
	)

	err = row.Scan(&deptID, &name, &description, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find department by id: %w", err)
	}

	return r.mapToDepartment(deptID, name, description, createdAt, updatedAt)
}

func (r *DepartmentRepository) FindAll(ctx context.Context) ([]*department.Department, error) {
	query, args, err := r.builder.Select("id", "name", "description", "created_at", "updated_at").
		From("departments").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build find all departments query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("find all departments: %w", err)
	}
	defer rows.Close()

	var departments []*department.Department
	for rows.Next() {
		var (
			deptID      string
			name        string
			description string
			createdAt   time.Time
			updatedAt   time.Time
		)

		if err := rows.Scan(&deptID, &name, &description, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan department row: %w", err)
		}

		d, err := r.mapToDepartment(deptID, name, description, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate department rows: %w", err)
	}

	return departments, nil
}

func (r *DepartmentRepository) Delete(ctx context.Context, id department.DepartmentID) error {
	query, args, err := r.builder.Delete("departments").
		Where(sq.Eq{"id": id.UUID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete department query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete department: %w", err)
	}

	return nil
}

func (r *DepartmentRepository) mapToDepartment(
	deptID, name, description string,
	createdAt, updatedAt time.Time,
) (*department.Department, error) {
	id, err := department.ParseDepartmentID(deptID)
	if err != nil {
		return nil, err
	}

	n, err := department.NewDepartmentName(name)
	if err != nil {
		return nil, err
	}

	desc, err := department.NewDepartmentDescription(description)
	if err != nil {
		return nil, err
	}

	return department.ReconstituteDepartment(id, n, desc, createdAt, updatedAt), nil
}
