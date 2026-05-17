package department

import (
	"time"
)

type Department struct {
	id          DepartmentID
	name        DepartmentName
	description DepartmentDescription
	createdAt   time.Time
	updatedAt   time.Time
	events      []DomainEvent
}

func NewDepartment(name DepartmentName, description DepartmentDescription) *Department {
	now := time.Now()
	id := NewDepartmentID()
	d := &Department{
		id:        id,
		name:      name,
		createdAt: now,
		updatedAt: now,
	}
	if description.String() != "" {
		d.description = description
	}
	d.emit(NewDepartmentCreated(id, name.String()))
	return d
}

func ReconstituteDepartment(
	id DepartmentID,
	name DepartmentName,
	description DepartmentDescription,
	createdAt time.Time,
	updatedAt time.Time,
) *Department {
	return &Department{
		id:          id,
		name:        name,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (d *Department) ID() DepartmentID                    { return d.id }
func (d *Department) Name() DepartmentName                 { return d.name }
func (d *Department) Description() DepartmentDescription   { return d.description }
func (d *Department) CreatedAt() time.Time                 { return d.createdAt }
func (d *Department) UpdatedAt() time.Time                 { return d.updatedAt }

func (d *Department) UpdateName(name DepartmentName) {
	old := d.name.String()
	d.name = name
	d.updatedAt = time.Now()
	d.emit(DepartmentNameChanged{
		baseEvent:    baseEvent{occurredAt: time.Now()},
		DepartmentID: d.id,
		OldName:      old,
		NewName:      name.String(),
	})
}

func (d *Department) UpdateDescription(desc DepartmentDescription) {
	old := d.description.String()
	d.description = desc
	d.updatedAt = time.Now()
	d.emit(DepartmentDescriptionChanged{
		baseEvent:      baseEvent{occurredAt: time.Now()},
		DepartmentID:   d.id,
		OldDescription: old,
		NewDescription: desc.String(),
	})
}

func (d *Department) Delete() {
	d.emit(NewDepartmentDeleted(d.id))
}

func (d *Department) emit(event DomainEvent) {
	d.events = append(d.events, event)
}

func (d *Department) PullEvents() []DomainEvent {
	events := d.events
	d.events = nil
	return events
}
