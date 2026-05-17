package taskhttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"github.com/wattanar/taskmanager/api/spec"
	taskapp "github.com/wattanar/taskmanager/internal/application/task"
	domainErrors "github.com/wattanar/taskmanager/internal/domain/errors"
	"github.com/wattanar/taskmanager/internal/domain/task"
)

type TaskHandler struct {
	createUC *taskapp.CreateTaskUseCase
	getUC    *taskapp.GetTaskUseCase
	listUC   *taskapp.ListTasksUseCase
	updateUC *taskapp.UpdateTaskUseCase
	deleteUC *taskapp.DeleteTaskUseCase
}

func NewTaskHandler(
	createUC *taskapp.CreateTaskUseCase,
	getUC *taskapp.GetTaskUseCase,
	listUC *taskapp.ListTasksUseCase,
	updateUC *taskapp.UpdateTaskUseCase,
	deleteUC *taskapp.DeleteTaskUseCase,
) *TaskHandler {
	return &TaskHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req spec.CreateTaskJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	priority := task.PriorityMedium
	if req.Priority != nil {
		priority = task.Priority(*req.Priority)
	}

	t, err := h.createUC.Execute(r.Context(), taskapp.CreateTaskCommand{
		Title:       req.Title,
		Description: fromPtr(req.Description),
		Priority:    priority,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, mapToSpecTask(t))
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request, id types.UUID) {
	taskID, err := task.ParseTaskID(id.String())
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	t, err := h.getUC.Execute(r.Context(), taskapp.GetTaskQuery{ID: taskID})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapToSpecTask(t))
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request, params spec.ListTasksParams) {
	q := taskapp.ListTasksQuery{}

	if params.Status != nil {
		s := task.TaskStatus(*params.Status)
		q.Status = &s
	}

	if params.Priority != nil {
		p := task.Priority(*params.Priority)
		q.Priority = &p
	}

	tasks, err := h.listUC.Execute(r.Context(), q)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	result := make([]spec.Task, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, mapToSpecTask(t))
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request, id types.UUID) {
	taskID, err := task.ParseTaskID(id.String())
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req spec.UpdateTaskJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	cmd := taskapp.UpdateTaskCommand{ID: taskID}

	if req.Title != nil {
		cmd.Title = req.Title
	}

	if req.Description != nil {
		cmd.Description = req.Description
	}

	if req.Status != nil {
		s := task.TaskStatus(*req.Status)
		cmd.Status = &s
	}

	if req.Priority != nil {
		p := task.Priority(*req.Priority)
		cmd.Priority = &p
	}

	t, err := h.updateUC.Execute(r.Context(), cmd)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapToSpecTask(t))
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request, id types.UUID) {
	taskID, err := task.ParseTaskID(id.String())
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	if err := h.deleteUC.Execute(r.Context(), taskapp.DeleteTaskCommand{ID: taskID}); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, spec.Error{Message: message})
}

func writeDomainError(w http.ResponseWriter, err error) {
	var notFound *domainErrors.NotFound
	var invalidArg *domainErrors.InvalidArgument
	var invalidTransition *domainErrors.InvalidTransition

	if errors.As(err, &notFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if errors.As(err, &invalidArg) || errors.As(err, &invalidTransition) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeError(w, http.StatusInternalServerError, "internal server error")
}

func mapToSpecTask(t *task.Task) spec.Task {
	id, _ := uuid.Parse(t.ID().String())

	desc := t.Description().String()

	return spec.Task{
		Id:          id,
		Title:       t.Title().String(),
		Description: &desc,
		Status:      spec.TaskStatus(t.Status()),
		Priority:    spec.Priority(t.Priority()),
		CreatedAt:   t.CreatedAt(),
		UpdatedAt:   t.UpdatedAt(),
	}
}

func fromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
