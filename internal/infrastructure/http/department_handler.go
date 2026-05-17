package taskhttp

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"github.com/wattanar/taskmanager/api/spec"
	deptapp "github.com/wattanar/taskmanager/internal/application/department"
	"github.com/wattanar/taskmanager/internal/domain/department"
)

type DepartmentHandler struct {
	createUC *deptapp.CreateDepartmentUseCase
	getUC    *deptapp.GetDepartmentUseCase
	listUC   *deptapp.ListDepartmentsUseCase
	updateUC *deptapp.UpdateDepartmentUseCase
	deleteUC *deptapp.DeleteDepartmentUseCase
}

func NewDepartmentHandler(
	createUC *deptapp.CreateDepartmentUseCase,
	getUC *deptapp.GetDepartmentUseCase,
	listUC *deptapp.ListDepartmentsUseCase,
	updateUC *deptapp.UpdateDepartmentUseCase,
	deleteUC *deptapp.DeleteDepartmentUseCase,
) *DepartmentHandler {
	return &DepartmentHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req spec.CreateDepartmentJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	d, err := h.createUC.Execute(r.Context(), deptapp.CreateDepartmentCommand{
		Name:        req.Name,
		Description: fromPtr(req.Description),
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, mapToSpecDepartment(d))
}

func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request, id types.UUID) {
	deptID, err := department.ParseDepartmentID(id.String())
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	d, err := h.getUC.Execute(r.Context(), deptapp.GetDepartmentQuery{ID: deptID})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapToSpecDepartment(d))
}

func (h *DepartmentHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	departments, err := h.listUC.Execute(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}

	result := make([]spec.Department, 0, len(departments))
	for _, d := range departments {
		result = append(result, mapToSpecDepartment(d))
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request, id types.UUID) {
	deptID, err := department.ParseDepartmentID(id.String())
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	var req spec.UpdateDepartmentJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	cmd := deptapp.UpdateDepartmentCommand{ID: deptID}

	if req.Name != nil {
		cmd.Name = req.Name
	}

	if req.Description != nil {
		cmd.Description = req.Description
	}

	d, err := h.updateUC.Execute(r.Context(), cmd)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapToSpecDepartment(d))
}

func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request, id types.UUID) {
	deptID, err := department.ParseDepartmentID(id.String())
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	if err := h.deleteUC.Execute(r.Context(), deptapp.DeleteDepartmentCommand{ID: deptID}); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func mapToSpecDepartment(d *department.Department) spec.Department {
	id, _ := uuid.Parse(d.ID().String())

	desc := d.Description().String()

	return spec.Department{
		Id:          id,
		Name:        d.Name().String(),
		Description: &desc,
		CreatedAt:   d.CreatedAt(),
		UpdatedAt:   d.UpdatedAt(),
	}
}
