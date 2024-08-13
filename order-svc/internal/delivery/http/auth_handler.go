package http

import (
	"order-svc/internal/dto"
	"order-svc/pkg"

	"github.com/benebobaa/valo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (ah *AuthHandler) Authenticate(c *gin.Context) {

	var req dto.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := valo.Validate(req)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userInfo := pkg.UserInfo{
		ID:       uuid.New().String(),
		Username: req.Username,
	}
	token, err := pkg.GenerateToken(userInfo)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dto.AuthResponse{
		User:  userInfo,
		Token: token,
	})
}
