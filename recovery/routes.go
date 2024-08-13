package http

import (
	"order-svc/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.Use(middleware.AuthMiddleware())
	group.POST("/order", h.CreateOrder)
}

func (a *AuthHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/guest", a.Authenticate)
}
