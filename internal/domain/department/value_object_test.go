package department_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainErrors "github.com/wattanar/taskmanager/internal/domain/errors"
	"github.com/wattanar/taskmanager/internal/domain/department"
)

func TestNewDepartmentName(t *testing.T) {
	t.Run("valid name", func(t *testing.T) {
		name, err := department.NewDepartmentName("Engineering")
		require.NoError(t, err)
		assert.Equal(t, "Engineering", name.String())
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := department.NewDepartmentName("")
		require.Error(t, err)
		var invalidArg *domainErrors.InvalidArgument
		assert.ErrorAs(t, err, &invalidArg)
	})

	t.Run("name too long", func(t *testing.T) {
		long := string(make([]byte, 201))
		_, err := department.NewDepartmentName(long)
		require.Error(t, err)
		var invalidArg *domainErrors.InvalidArgument
		assert.ErrorAs(t, err, &invalidArg)
	})
}

func TestNewDepartmentDescription(t *testing.T) {
	t.Run("valid description", func(t *testing.T) {
		desc, err := department.NewDepartmentDescription("Handles all engineering tasks")
		require.NoError(t, err)
		assert.Equal(t, "Handles all engineering tasks", desc.String())
	})

	t.Run("empty description", func(t *testing.T) {
		desc, err := department.NewDepartmentDescription("")
		require.NoError(t, err)
		assert.Equal(t, "", desc.String())
	})

	t.Run("description too long", func(t *testing.T) {
		long := string(make([]byte, 2001))
		_, err := department.NewDepartmentDescription(long)
		require.Error(t, err)
	})
}

func TestParseDepartmentID(t *testing.T) {
	t.Run("valid uuid", func(t *testing.T) {
		id, err := department.ParseDepartmentID("550e8400-e29b-41d4-a716-446655440000")
		require.NoError(t, err)
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
	})

	t.Run("invalid uuid", func(t *testing.T) {
		_, err := department.ParseDepartmentID("not-a-uuid")
		require.Error(t, err)
	})
}
