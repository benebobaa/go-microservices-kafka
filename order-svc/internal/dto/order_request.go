package dto

type Status int

const (
	PENDING Status = iota
	PROCESSING
	CANCEL_PROCESSING
	COMPLETE
	CANCELLED
)

func (s Status) String() string {
	return [...]string{"PENDING", "PROCESSING", "CANCEL_PROCESSING", "COMPLETE", "CANCELLED"}[s]
}

type OrderRequest struct {
	ProductID  string `json:"product_id" valo:"notblank"`
	Quantity   int32  `json:"quantity" valo:"min=1"`
	CustomerID string `json:"-"`
	Username   string `json:"-"`
}

type OrderUpdateRequest struct {
	RefID     string  `json:"ref_id"`
	Amount    float64 `json:"amount"`
	Quantity  int32   `json:"quantity"`
	Status    string  `json:"-"`
	EventType string  `json:"-"`
}

type OrderCancelRequest struct {
	OrderID  int    `json:"order_id" valo:"min=1"`
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
