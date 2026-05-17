package department

import "context"

type DepartmentRepository interface {
	Save(ctx context.Context, department *Department) error
	FindByID(ctx context.Context, id DepartmentID) (*Department, error)
	FindAll(ctx context.Context) ([]*Department, error)
	Delete(ctx context.Context, id DepartmentID) error
}
