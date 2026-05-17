package department_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wattanar/taskmanager/internal/domain/department"
)

func TestNewDepartment(t *testing.T) {
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("Engineering department")
	d := department.NewDepartment(name, desc)

	assert.Equal(t, "Engineering", d.Name().String())
	assert.Equal(t, "Engineering department", d.Description().String())
	assert.False(t, d.CreatedAt().IsZero())
	assert.False(t, d.UpdatedAt().IsZero())

	events := d.PullEvents()
	assert.Len(t, events, 1)
	assert.Equal(t, "department.created", events[0].EventName())
}

func TestDepartment_UpdateName(t *testing.T) {
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("")
	d := department.NewDepartment(name, desc)
	d.PullEvents()

	newName, _ := department.NewDepartmentName("Engineering II")
	d.UpdateName(newName)

	assert.Equal(t, "Engineering II", d.Name().String())

	events := d.PullEvents()
	assert.Len(t, events, 1)
	assert.Equal(t, "department.name_changed", events[0].EventName())
}

func TestDepartment_UpdateDescription(t *testing.T) {
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("Original")
	d := department.NewDepartment(name, desc)
	d.PullEvents()

	newDesc, _ := department.NewDepartmentDescription("Updated description")
	d.UpdateDescription(newDesc)

	assert.Equal(t, "Updated description", d.Description().String())
}

func TestDepartment_Timestamps(t *testing.T) {
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("")
	d := department.NewDepartment(name, desc)

	originalUpdated := d.UpdatedAt()
	time.Sleep(time.Millisecond)

	newName, _ := department.NewDepartmentName("Engineering II")
	d.UpdateName(newName)

	assert.True(t, d.UpdatedAt().After(originalUpdated))
}

func TestDepartment_PullEvents(t *testing.T) {
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("")
	d := department.NewDepartment(name, desc)

	events := d.PullEvents()
	assert.Len(t, events, 1)

	empty := d.PullEvents()
	assert.Len(t, empty, 0)
}

func TestReconstituteDepartment(t *testing.T) {
	now := time.Now()
	id := department.NewDepartmentID()
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("Desc")

	d := department.ReconstituteDepartment(id, name, desc, now, now)

	assert.Equal(t, id, d.ID())
	assert.Equal(t, "Engineering", d.Name().String())
	assert.Equal(t, "Desc", d.Description().String())
	assert.Equal(t, now, d.CreatedAt())
	assert.Equal(t, now, d.UpdatedAt())

	events := d.PullEvents()
	assert.Len(t, events, 0)
}
