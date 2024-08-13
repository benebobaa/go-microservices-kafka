package http

import (
	"user-svc/internal/dto"
	"user-svc/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase *usecase.Usecase
}

func NewHandler(usecase *usecase.Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) TestValidate(c *gin.Context) {

	username := c.DefaultQuery("username", "beneboba")

	response, err := h.usecase.ValidateUser(
		c,
		&dto.UserValidateRequest{Username: username},
	)

	if err != nil {
		c.JSON(400, response)
		return
	}

	c.JSON(200, response)
}
