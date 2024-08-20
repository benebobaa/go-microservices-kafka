package http

import (
	"github.com/benebobaa/valo"
	"github.com/gin-gonic/gin"
	"order-svc/internal/dto"
	"order-svc/internal/usecase"
)

type BankRegisHandler struct {
	usecase *usecase.BankRegistrationUsecase
}

func NewBankRegisHandler(usecase *usecase.BankRegistrationUsecase) *BankRegisHandler {
	return &BankRegisHandler{usecase: usecase}
}

func (brh *BankRegisHandler) RegisterBankAccount(c *gin.Context) {

	var req dto.BankRegistrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := valo.Validate(req)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, err := brh.usecase.RegisterBankAccount(c, &req)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}
