package event

import (
	"database/sql"
	"time"
)

type Order struct {
	ID          int32          `json:"id"`
	CustomerID  string         `json:"customer_id"`
	Username    string         `json:"username"`
	ProductName string         `json:"product_name"`
	OrderDate   time.Time      `json:"order_date"`
	Status      string         `json:"status"`
	TotalAmount sql.NullString `json:"total_amount"`
}

type OrderRequest struct {
	ProductID  string `json:"product_id"`
	Quantity   int32  `json:"quantity"`
	CustomerID string `json:"-"`
	Username   string `json:"-"`
}
