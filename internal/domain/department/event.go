package department

import "time"

type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

type baseEvent struct {
	occurredAt time.Time
}

func (e baseEvent) OccurredAt() time.Time { return e.occurredAt }

type DepartmentCreated struct {
	baseEvent
	DepartmentID DepartmentID
	Name         string
}

func NewDepartmentCreated(id DepartmentID, name string) DepartmentCreated {
	return DepartmentCreated{
		baseEvent:    baseEvent{occurredAt: time.Now()},
		DepartmentID: id,
		Name:         name,
	}
}

func (e DepartmentCreated) EventName() string { return "department.created" }

type DepartmentNameChanged struct {
	baseEvent
	DepartmentID DepartmentID
	OldName      string
	NewName      string
}

func (e DepartmentNameChanged) EventName() string { return "department.name_changed" }

type DepartmentDescriptionChanged struct {
	baseEvent
	DepartmentID DepartmentID
	OldDescription string
	NewDescription string
}

func (e DepartmentDescriptionChanged) EventName() string { return "department.description_changed" }

type DepartmentDeleted struct {
	baseEvent
	DepartmentID DepartmentID
}

func NewDepartmentDeleted(id DepartmentID) DepartmentDeleted {
	return DepartmentDeleted{
		baseEvent:    baseEvent{occurredAt: time.Now()},
		DepartmentID: id,
	}
}

func (e DepartmentDeleted) EventName() string { return "department.deleted" }
