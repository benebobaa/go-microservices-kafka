package handler

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/benebobaa/valo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	transactionFile = "transactions.json"
	balanceFile     = "balances.json"
)

type Transaction struct {
	ID            string  `json:"id"`
	RefID         string  `json:"ref_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	AccountBankID string  `json:"account_bank_id"`
}

type Balance struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type TransactionRequest struct {
	RefID         string  `json:"ref_id" valo:"notblank,sizeMin=3"`
	Amount        float64 `json:"amount" valo:"min=1"`
	AccountBankID string  `json:"account_bank_id" valo:"notblank,sizeMin=3"`
}

type RefundRequest struct {
	RefID string `json:"ref_id" valo:"notblank,sizeMin=3"`
}

type PaymentHandler struct {
	dbT map[string]Transaction
	dbB map[string]Balance

	mutex *sync.RWMutex
}

func (h *PaymentHandler) GetBalance(c *gin.Context) {
	var balances []Balance

	for _, balance := range h.dbB {
		balances = append(balances, balance)
	}

	c.JSON(200, gin.H{"status_code": 200, "data": balances})
}

func (h *PaymentHandler) GetTransaction(c *gin.Context) {
	var transactions []Transaction

	for _, v := range h.dbT {
		transactions = append(transactions, v)
	}

	c.JSON(200, gin.H{"status_code": 200, "data": transactions})
}

func (h *PaymentHandler) CreateTransaction(c *gin.Context) {
	var req TransactionRequest

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

	for _, transaction := range h.dbT {
		if transaction.RefID == req.RefID {
			c.JSON(400, gin.H{"status_code": 400, "error": "ref_id already exists"})
			return
		}
	}

	account, ok := h.dbB[req.AccountBankID]

	if !ok {
		c.JSON(404, gin.H{"status_code": 404, "error": "account not found"})
		return
	}

	if account.Balance < req.Amount {
		c.JSON(400, gin.H{"status_code": 400, "error": "balance is not enough"})
		return
	}

	account.Balance -= req.Amount
	h.dbB[account.AccountID] = account

	transaction := Transaction{
		ID:            uuid.New().String(),
		RefID:         req.RefID,
		Amount:        req.Amount,
		Status:        "success",
		AccountBankID: req.AccountBankID,
	}

	h.dbT[transaction.ID] = transaction

	c.JSON(201, gin.H{"status_code": 201, "data": transaction})
}

func (h *PaymentHandler) RefundTransaction(c *gin.Context) {
	var req RefundRequest

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

	for _, transaction := range h.dbT {
		if transaction.RefID == req.RefID {
			account, ok := h.dbB[transaction.AccountBankID]

			if !ok {
				c.JSON(404, gin.H{"status_code": 404, "error": "account not found"})
				return
			}

			if transaction.Status == "refunded" {
				c.JSON(400, gin.H{"status_code": 400, "error": "transaction already refunded"})
				return
			}

			account.Balance += transaction.Amount
			h.dbB[transaction.AccountBankID] = account

			transaction.Status = "refunded"
			h.dbT[transaction.ID] = transaction

			c.JSON(200, gin.H{"status_code": 200, "data": transaction})
			return
		}
	}

	c.JSON(404, gin.H{"status_code": 404, "error": "transaction not found"})
}

// Add these functions to the PaymentHandler struct
func (h *PaymentHandler) SaveData() error {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// Save transactions
	if err := saveToFile(transactionFile, h.dbT); err != nil {
		return err
	}

	// Save balances
	if err := saveToFile(balanceFile, h.dbB); err != nil {
		return err
	}

	return nil
}

func (h *PaymentHandler) LoadData() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Load transactions
	if err := loadFromFile(transactionFile, &h.dbT); err != nil {
		return err
	}

	// Load balances
	if err := loadFromFile(balanceFile, &h.dbB); err != nil {
		return err
	}

	balance1 := Balance{
		AccountID: "AC-001",
		Balance:   50000,
	}

	balance2 := Balance{
		AccountID: "AC-002",
		Balance:   100000,
	}

	balance3 := Balance{
		AccountID: "AC-003",
		Balance:   5000,
	}

	var balances []Balance = []Balance{balance1, balance2, balance3}

	for _, balance := range balances {
		_, exists := h.dbB[balance.AccountID]
		if !exists {
			h.dbB[balance.AccountID] = balance
		}
	}

	return nil
}

func saveToFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func loadFromFile(filename string, data interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file to load, start with an empty map
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(data)
}
