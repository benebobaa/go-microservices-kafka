package handler

import (
	"sync"

	"github.com/benebobaa/valo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	ID            string `json:"id"`
	AccountBankID string `json:"account_bank_id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
}

type UserRequest struct {
	Username string `json:"username" valo:"notblank,sizeMin=3"`
	Email    string `json:"email" valo:"notblank,email"`
}

type UserUpdateRequest struct {
	AccountBankID string `json:"account_bank_id" valo:"notblank"`
}

type UserHandler struct {
	db    map[string]User
	mutex *sync.RWMutex
}

//
// func (h *UserHandler) UpdateUserBankID(c *gin.Context) {
//
//
// }

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status_code": 400, "error": err.Error()})
		return
	}

	err := valo.Validate(req)

	if err != nil {
		c.JSON(400, gin.H{"status_code": 400, "error": err.Error()})
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, user := range h.db {
		if user.Username == req.Username {
			c.JSON(400, gin.H{"status_code": 400, "error": "username already exists"})
			return
		}

		if user.Email == req.Email {
			c.JSON(400, gin.H{"status_code": 400, "error": "email already exists"})
			return
		}
	}

	user := User{
		ID:            "USER-" + uuid.New().String(),
		AccountBankID: "",
		Username:      req.Username,
		Email:         req.Email,
	}

	h.db[user.ID] = user

	c.JSON(201, gin.H{"status_code": 201, "data": user})
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
