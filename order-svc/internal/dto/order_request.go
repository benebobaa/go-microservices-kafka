package dto

import "database/sql"

type Status int

const (
	PENDING Status = iota
	PROCESSING
	CANCEL_PROCESSING
	COMPLETE
	FAILED
)

func (s Status) String() string {
	return [...]string{"PENDING", "PROCESSING", "CANCEL_PROCESSING", "COMPLETE", "FAILED"}[s]
}

type OrderRequest struct {
	ProductID  string `json:"product_id" valo:"notblank"`
	Quantity   int32  `json:"quantity" valo:"min=1"`
	CustomerID string `json:"-"`
	Username   string `json:"-"`
}

type OrderUpdateRequest struct {
	OrderRefID  string
	Status      string          `json:"-"`
	TotalAmount sql.NullFloat64 `json:"total_amount"`
}

type OrderCancelRequest struct {
	OrderID  int32  `json:"order_id" valo:"min=1"`
	Username string `json:"-"`
}

// Entity
// type Order struct {
// 	ID          int32          `json:"id"`
// 	CustomerID  string         `json:"customer_id"`
// 	Username    string         `json:"username"`
// 	ProductName string         `json:"product_name"`
// 	OrderDate   time.Time      `json:"order_date"`
// 	Status      string         `json:"status"`
// 	TotalAmount sql.NullString `json:"total_amount"`
// }
