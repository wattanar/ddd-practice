package taskhttp_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wattanar/taskmanager/api/spec"
	deptapp "github.com/wattanar/taskmanager/internal/application/department"
	taskapp "github.com/wattanar/taskmanager/internal/application/task"
	"github.com/wattanar/taskmanager/internal/domain/department"
	"github.com/wattanar/taskmanager/internal/domain/department/mocks"
	"github.com/wattanar/taskmanager/internal/domain/task"
	taskmocks "github.com/wattanar/taskmanager/internal/domain/task/mocks"
	taskhttp "github.com/wattanar/taskmanager/internal/infrastructure/http"
)

func TestCreateDepartmentHandler(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)
	mockRepo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*department.Department")).Return(nil).Once()

	handler := newDeptHandler(t, mockRepo)

	body := `{"name":"Engineering"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/departments", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp spec.Department
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "Engineering", resp.Name)
}

func TestCreateDepartmentHandler_EmptyName(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)
	handler := newDeptHandler(t, mockRepo)

	body := `{"name":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/departments", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp spec.Error
	json.NewDecoder(w.Body).Decode(&errResp)
	assert.Contains(t, errResp.Message, "invalid argument")
}

func TestGetDepartmentHandler(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)
	id := uuid.New()
	domainID, _ := department.ParseDepartmentID(id.String())
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("Desc")
	expected := department.ReconstituteDepartment(domainID, name, desc, now, now)

	mockRepo.EXPECT().FindByID(mock.Anything, domainID).Return(expected, nil).Once()

	handler := newDeptHandler(t, mockRepo)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/"+id.String(), nil)
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp spec.Department
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Engineering", resp.Name)
}

func TestGetDepartmentHandler_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)
	id := uuid.New()
	domainID, _ := department.ParseDepartmentID(id.String())

	mockRepo.EXPECT().FindByID(mock.Anything, domainID).Return(nil, nil).Once()

	handler := newDeptHandler(t, mockRepo)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/"+id.String(), nil)
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListDepartmentsHandler(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)
	mockRepo.EXPECT().FindAll(mock.Anything).Return([]*department.Department{}, nil).Once()

	handler := newDeptHandler(t, mockRepo)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments", nil)
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var departments []spec.Department
	json.NewDecoder(w.Body).Decode(&departments)
	assert.Empty(t, departments)
}

func TestUpdateDepartmentHandler(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)
	id := uuid.New()
	domainID, _ := department.ParseDepartmentID(id.String())
	name, _ := department.NewDepartmentName("Old")
	desc, _ := department.NewDepartmentDescription("")
	existing := department.ReconstituteDepartment(domainID, name, desc, now, now)

	mockRepo.EXPECT().FindByID(mock.Anything, domainID).Return(existing, nil).Once()
	mockRepo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*department.Department")).Return(nil).Once()

	handler := newDeptHandler(t, mockRepo)

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/departments/"+id.String(),
		bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp spec.Department
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "New Name", resp.Name)
}

func TestDeleteDepartmentHandler(t *testing.T) {
	mockRepo := mocks.NewMockDepartmentRepository(t)
	id := uuid.New()
	domainID, _ := department.ParseDepartmentID(id.String())
	name, _ := department.NewDepartmentName("Engineering")
	desc, _ := department.NewDepartmentDescription("")
	existing := department.ReconstituteDepartment(domainID, name, desc, now, now)

	mockRepo.EXPECT().FindByID(mock.Anything, domainID).Return(existing, nil).Once()
	mockRepo.EXPECT().Delete(mock.Anything, domainID).Return(nil).Once()

	handler := newDeptHandler(t, mockRepo)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/departments/"+id.String(), nil)
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func newDeptHandler(t *testing.T, repo department.DepartmentRepository) *taskhttp.APIHandler {
	deptHandler := taskhttp.NewDepartmentHandler(
		deptapp.NewCreateDepartmentUseCase(repo),
		deptapp.NewGetDepartmentUseCase(repo),
		deptapp.NewListDepartmentsUseCase(repo),
		deptapp.NewUpdateDepartmentUseCase(repo),
		deptapp.NewDeleteDepartmentUseCase(repo),
	)

	taskMock := taskmocks.NewMockTaskRepository(t)
	taskHandler := taskhttp.NewTaskHandler(
		taskapp.NewCreateTaskUseCase(taskMock),
		taskapp.NewGetTaskUseCase(taskMock),
		taskapp.NewListTasksUseCase(taskMock),
		taskapp.NewUpdateTaskUseCase(taskMock),
		taskapp.NewDeleteTaskUseCase(taskMock),
	)

	return taskhttp.NewAPIHandler(taskHandler, deptHandler)
}

var _ department.DepartmentRepository = (*mocks.MockDepartmentRepository)(nil)
var _ task.TaskRepository = (*taskmocks.MockTaskRepository)(nil)
