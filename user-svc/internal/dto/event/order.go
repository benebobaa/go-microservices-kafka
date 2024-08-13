package event

import (
	"database/sql"
	"time"
)

type Order struct {
	ID          int32           `json:"id"`
	OrderRefID  string          `json:"order_ref_id"`
	CustomerID  string          `json:"customer_id"`
	Username    string          `json:"username"`
	ProductID   string          `json:"product_id"`
	Quantity    int32           `json:"quantity"`
	OrderDate   time.Time       `json:"order_date"`
	Status      string          `json:"status"`
	TotalAmount sql.NullFloat64 `json:"total_amount"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
