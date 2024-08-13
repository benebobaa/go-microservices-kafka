package http

import (
	"orchestra-svc/internal/dto"
	"orchestra-svc/internal/usecase"

	"github.com/benebobaa/valo"
	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	wu *usecase.WorkflowUsecase
}

func NewWorkflowHandler(wu *usecase.WorkflowUsecase) *WorkflowHandler {
	return &WorkflowHandler{wu: wu}
}

func (wf *WorkflowHandler) CreateWorkflow(c *gin.Context) {

	var req dto.WorkflowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := valo.Validate(req)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, err := wf.wu.CreateWorkflow(c, &req)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, response)
}

func (wf *WorkflowHandler) GetStepsByType(c *gin.Context) {

	name := c.Query("type")

	if name == "" {
		c.JSON(400, gin.H{"error": "type is required"})
		return
	}

	response, err := wf.wu.GetByType(c, name)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}
