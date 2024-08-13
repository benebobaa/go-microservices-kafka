package http

import "github.com/gin-gonic/gin"

func (wh *WorkflowHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/", wh.CreateWorkflow)
	router.GET("/steps", wh.GetStepsByType)
}
