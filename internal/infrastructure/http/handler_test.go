package taskhttp_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

var now = time.Now().UTC()

func TestCreateTaskHandler(t *testing.T) {
	mockRepo := taskmocks.NewMockTaskRepository(t)
	mockRepo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*task.Task")).Return(nil).Once()

	handler := newHandler(t, mockRepo)

	body := `{"title":"Buy groceries","priority":"high"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp spec.Task
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "Buy groceries", resp.Title)
	assert.Equal(t, spec.Todo, resp.Status)
	assert.Equal(t, spec.High, resp.Priority)
}

func TestCreateTaskHandler_EmptyTitle(t *testing.T) {
	mockRepo := taskmocks.NewMockTaskRepository(t)
	handler := newHandler(t, mockRepo)

	body := `{"title":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp spec.Error
	json.NewDecoder(w.Body).Decode(&errResp)
	assert.Contains(t, errResp.Message, "invalid argument")
}

func TestGetTaskHandler(t *testing.T) {
	mockRepo := taskmocks.NewMockTaskRepository(t)
	id := uuid.New()
	domainID, _ := task.ParseTaskID(id.String())
	title, _ := task.NewTitle("Test")
	desc, _ := task.NewDescription("Desc")
	expected := task.ReconstituteTask(domainID, title, desc, task.TaskStatusInProgress, task.PriorityMedium, now, now)

	mockRepo.EXPECT().FindByID(mock.Anything, domainID).Return(expected, nil).Once()

	handler := newHandler(t, mockRepo)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+id.String(), nil)
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp spec.Task
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Test", resp.Title)
	assert.Equal(t, spec.InProgress, resp.Status)
}

func TestGetTaskHandler_NotFound(t *testing.T) {
	mockRepo := taskmocks.NewMockTaskRepository(t)
	id := uuid.New()
	domainID, _ := task.ParseTaskID(id.String())

	mockRepo.EXPECT().FindByID(mock.Anything, domainID).Return(nil, nil).Once()

	handler := newHandler(t, mockRepo)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+id.String(), nil)
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListTasksHandler(t *testing.T) {
	mockRepo := taskmocks.NewMockTaskRepository(t)
	mockRepo.EXPECT().FindAll(mock.Anything, mock.AnythingOfType("task.TaskFilter")).Return([]*task.Task{}, nil).Once()

	handler := newHandler(t, mockRepo)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var tasks []spec.Task
	json.NewDecoder(w.Body).Decode(&tasks)
	assert.Empty(t, tasks)
}

func TestUpdateTaskHandler(t *testing.T) {
	mockRepo := taskmocks.NewMockTaskRepository(t)
	id := uuid.New()
	domainID, _ := task.ParseTaskID(id.String())
	title, _ := task.NewTitle("Old")
	desc, _ := task.NewDescription("")
	existing := task.ReconstituteTask(domainID, title, desc, task.TaskStatusTodo, task.PriorityLow, now, now)

	mockRepo.EXPECT().FindByID(mock.Anything, domainID).Return(existing, nil).Once()
	mockRepo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*task.Task")).Return(nil).Once()

	handler := newHandler(t, mockRepo)

	newStatus := "in_progress"
	body := `{"status":"` + newStatus + `","priority":"critical"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/"+id.String(),
		bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp spec.Task
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, spec.InProgress, resp.Status)
	assert.Equal(t, spec.Critical, resp.Priority)
}

func TestUpdateTaskHandler_NotFound(t *testing.T) {
	mockRepo := taskmocks.NewMockTaskRepository(t)
	id := uuid.New()
	domainID, _ := task.ParseTaskID(id.String())

	mockRepo.EXPECT().FindByID(mock.Anything, domainID).Return(nil, nil).Once()

	handler := newHandler(t, mockRepo)
	body := `{"title":"New title"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/"+id.String(),
		bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteTaskHandler(t *testing.T) {
	mockRepo := taskmocks.NewMockTaskRepository(t)
	id := uuid.New()
	domainID, _ := task.ParseTaskID(id.String())
	title, _ := task.NewTitle("Test")
	desc, _ := task.NewDescription("")
	existing := task.ReconstituteTask(domainID, title, desc, task.TaskStatusTodo, task.PriorityMedium, now, now)

	mockRepo.EXPECT().FindByID(mock.Anything, domainID).Return(existing, nil).Once()
	mockRepo.EXPECT().Delete(mock.Anything, domainID).Return(nil).Once()

	handler := newHandler(t, mockRepo)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+id.String(), nil)
	w := httptest.NewRecorder()

	mux := spec.HandlerFromMux(handler, http.NewServeMux())
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func newHandler(t *testing.T, repo task.TaskRepository) *taskhttp.APIHandler {
	createUC := taskapp.NewCreateTaskUseCase(repo)
	getUC := taskapp.NewGetTaskUseCase(repo)
	listUC := taskapp.NewListTasksUseCase(repo)
	updateUC := taskapp.NewUpdateTaskUseCase(repo)
	deleteUC := taskapp.NewDeleteTaskUseCase(repo)
	taskHandler := taskhttp.NewTaskHandler(createUC, getUC, listUC, updateUC, deleteUC)

	deptMock := mocks.NewMockDepartmentRepository(t)
	deptHandler := taskhttp.NewDepartmentHandler(
		deptapp.NewCreateDepartmentUseCase(deptMock),
		deptapp.NewGetDepartmentUseCase(deptMock),
		deptapp.NewListDepartmentsUseCase(deptMock),
		deptapp.NewUpdateDepartmentUseCase(deptMock),
		deptapp.NewDeleteDepartmentUseCase(deptMock),
	)

	return taskhttp.NewAPIHandler(taskHandler, deptHandler)
}

var _ task.TaskRepository = (*taskmocks.MockTaskRepository)(nil)
var _ department.DepartmentRepository = (*mocks.MockDepartmentRepository)(nil)
