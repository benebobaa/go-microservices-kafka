package http

import (
	"errors"
	"order-svc/internal/dto"
	"order-svc/internal/middleware"
	"order-svc/internal/usecase"
	"order-svc/pkg"

	"github.com/benebobaa/valo"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase *usecase.OrderUsecase
}

func NewOrderHandler(usecase *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: usecase}
}

func (oh *OrderHandler) CreateOrder(c *gin.Context) {

	user := c.MustGet(middleware.ClaimsKey).(*pkg.UserInfo)

	var req dto.OrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := valo.Validate(req)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req.CustomerID = user.ID
	req.Username = user.Username
	response, err := oh.usecase.CreateOrder(c, &req)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

func (oh *OrderHandler) CancelOrder(c *gin.Context) {

	user := c.MustGet(middleware.ClaimsKey).(*pkg.UserInfo)
	var req dto.OrderCancelRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := valo.Validate(req)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req.Username = user.Username
	response, err := oh.usecase.CancelOrder(c, &req)

	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorizeCancelOrder) {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecase.ErrCannotCancelOrder) {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}
