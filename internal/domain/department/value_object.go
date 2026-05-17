package department

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/wattanar/taskmanager/internal/domain/errors"
)

type DepartmentID struct {
	uuid.UUID
}

func NewDepartmentID() DepartmentID {
	return DepartmentID{UUID: uuid.New()}
}

func ParseDepartmentID(s string) (DepartmentID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return DepartmentID{}, fmt.Errorf("parse department id: %w", err)
	}
	return DepartmentID{UUID: id}, nil
}

func (id DepartmentID) String() string { return id.UUID.String() }

type DepartmentName struct {
	value string
}

func NewDepartmentName(s string) (DepartmentName, error) {
	if len(s) == 0 {
		return DepartmentName{}, &errors.InvalidArgument{Field: "name", Reason: "must not be empty"}
	}
	if len(s) > 200 {
		return DepartmentName{}, &errors.InvalidArgument{Field: "name", Reason: "must not exceed 200 characters"}
	}
	return DepartmentName{value: s}, nil
}

func (n DepartmentName) String() string { return n.value }

type DepartmentDescription struct {
	value string
}

func NewDepartmentDescription(s string) (DepartmentDescription, error) {
	if len(s) > 2000 {
		return DepartmentDescription{}, &errors.InvalidArgument{Field: "description", Reason: "must not exceed 2000 characters"}
	}
	return DepartmentDescription{value: s}, nil
}

func (d DepartmentDescription) String() string { return d.value }
