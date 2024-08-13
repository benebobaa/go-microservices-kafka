package http

import (
	"log"
	kafka "order-svc/internal/delivery/messaging"
	"order-svc/internal/dto"
	"order-svc/internal/dto/event"
	"order-svc/internal/middleware"
	"order-svc/internal/repository/sqlc"
	"order-svc/pkg"
	"time"

	"github.com/benebobaa/valo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	queries           sqlc.Querier
	orchestraProducer *kafka.KafkaProducer
}

func NewOrderHandler(queries *sqlc.Queries,
	orchestraProducer *kafka.KafkaProducer,
) *OrderHandler {
	return &OrderHandler{
		queries:           queries,
		orchestraProducer: orchestraProducer,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {

	var req dto.OrderRequest

	user := c.MustGet(middleware.ClaimsKey).(pkg.UserInfo)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := valo.Validate(req)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	order, err := h.queries.CreateOrder(c, sqlc.CreateOrderParams{
		ProductName: req.ProductName,
		Quantity:    int32(req.Quantity),
		CustomerID:  user.ID,
		Username:    user.Username,
	})

	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	event := event.GlobalEvent{
		EventID:   uuid.New().String(),
		EventType: "order_created",
		Timestamp: time.Now(),
		Source:    "order-service",
		Action:    "create",
		Status:    "success",
		Payload:   order,
	}

	eventBytes, err := event.ToJSON()
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.orchestraProducer.SendMessage(event.EventID, eventBytes)

	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(201, order)
}
