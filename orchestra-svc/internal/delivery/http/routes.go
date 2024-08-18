package http

import "github.com/gin-gonic/gin"

func (wf *WorkflowHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/", wf.CreateWorkflow)
	router.GET("/steps", wf.GetStepsByType)
	router.PATCH("/product-retry", wf.RetryProductReserve)
	router.PATCH("/retry", wf.RetryInstanceStep)
}
