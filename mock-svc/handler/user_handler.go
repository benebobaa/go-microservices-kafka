package handler

import "github.com/gin-gonic/gin"

type User struct {
	ID            string `json:"id"`
	AccountBankID string `json:"account_bank_id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
}

type UserHandler struct {
	db map[string]User
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var users []User

	for _, user := range h.db {
		users = append(users, user)
	}

	c.JSON(200, gin.H{"status_code": 200, "data": users})
}

func (h *UserHandler) GetuUserByUsername(c *gin.Context) {
	username := c.Param("username")

	for _, v := range h.db {
		if v.Username == username {
			c.JSON(200, gin.H{"status_code": 200, "data": v})
			return
		}
	}

	c.JSON(404, gin.H{"status_code": 404, "error": "user not found"})
}
