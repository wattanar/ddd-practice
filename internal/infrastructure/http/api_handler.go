package taskhttp

import (
	"github.com/wattanar/taskmanager/api/spec"
)

type APIHandler struct {
	*TaskHandler
	*DepartmentHandler
}

func NewAPIHandler(taskHandler *TaskHandler, deptHandler *DepartmentHandler) *APIHandler {
	return &APIHandler{
		TaskHandler:       taskHandler,
		DepartmentHandler: deptHandler,
	}
}

var _ spec.ServerInterface = (*APIHandler)(nil)
