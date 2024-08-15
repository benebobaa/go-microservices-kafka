package http

import (
	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
}

func NewWorkflowHandler() *WorkflowHandler {
	return &WorkflowHandler{}
}

func (wf *WorkflowHandler) CreateWorkflow(c *gin.Context) {

	c.JSON(200, "response")
}

func (wf *WorkflowHandler) GetStepsByType(c *gin.Context) {

	c.JSON(200, "response")
}
