package handler

import (
	"sync"

	"github.com/benebobaa/valo"
	"github.com/gin-gonic/gin"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
}

type ProductRequest struct {
	ProductID string `json:"product_id" valo:"notblank"`
	Quantity  int    `json:"quantity" valo:"min=1"`
}

type ProductHandler struct {
	db    map[string]Product
	mutex *sync.RWMutex
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	var products []Product

	for _, product := range h.db {
		products = append(products, product)
	}

	c.JSON(200, gin.H{"data": products})
}

func (h *ProductHandler) ReserveProduct(c *gin.Context) {
	var req ProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := valo.Validate(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.mutex.Lock()
	product, ok := h.db[req.ProductID]

	if !ok {
		c.JSON(404, gin.H{"error": "product not found"})
		return
	}

	if product.Stock < req.Quantity {
		c.JSON(400, gin.H{"error": "stock is not enough"})
		return
	}

	product.Stock -= req.Quantity

	h.db[req.ProductID] = product

	h.mutex.Unlock()

	c.JSON(200, gin.H{"data": product})
}

func (h *ProductHandler) ReleaseProduct(c *gin.Context) {
	var req ProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := valo.Validate(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.mutex.Lock()
	product, ok := h.db[req.ProductID]

	if !ok {
		c.JSON(404, gin.H{"error": "product not found"})
		return
	}

	product.Stock += req.Quantity
	h.db[req.ProductID] = product

	h.mutex.Unlock()

	c.JSON(200, gin.H{"data": product})
}
