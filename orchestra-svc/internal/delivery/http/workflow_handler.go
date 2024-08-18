package http

import (
	"github.com/gin-gonic/gin"
	"orchestra-svc/internal/dto"
	"orchestra-svc/internal/usecase"

	"github.com/benebobaa/valo"
)

type WorkflowHandler struct {
	oc *usecase.OrchestraUsecase
	rc *usecase.RetryUsecase
}

func NewWorkflowHandler(
	oc *usecase.OrchestraUsecase,
	rc *usecase.RetryUsecase,
) *WorkflowHandler {
	return &WorkflowHandler{
		oc: oc,
		rc: rc,
	}
}

func (wf *WorkflowHandler) CreateWorkflow(c *gin.Context) {

	c.JSON(200, "response")
}

func (wf *WorkflowHandler) GetStepsByType(c *gin.Context) {

	c.JSON(200, "response")
}

func (wf *WorkflowHandler) RetryProductReserve(c *gin.Context) {

	var req dto.ProductQuantityRetryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, err.Error())
		return
	}

	err := valo.Validate(req)

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	response, err := wf.rc.ProductQuantityRetry(c, &req)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, response)
}

func (wf *WorkflowHandler) RetryInstanceStep(c *gin.Context) {

	var req dto.RetryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, err.Error())
		return
	}

	err := valo.Validate(req)

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	response, err := wf.rc.RetryFailedInstanceStep(c, &req)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, response)
}
