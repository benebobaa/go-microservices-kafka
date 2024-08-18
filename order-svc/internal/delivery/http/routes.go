package http

import "github.com/gin-gonic/gin"

func (ah *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/guest", ah.Authenticate)
}

func (oh *OrderHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/", oh.CreateOrder)
	router.PATCH("/cancel", oh.CancelOrder)
	router.GET("/", oh.FindAll)
}
