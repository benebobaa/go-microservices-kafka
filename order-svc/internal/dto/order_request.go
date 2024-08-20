package dto

type Status int

const (
	PENDING Status = iota
	PROCESSING
	CANCEL_PROCESSING
	COMPLETE
	CANCELLED
	FAILED
)

func (s Status) String() string {
	return [...]string{"PENDING", "PROCESSING", "CANCEL_PROCESSING", "COMPLETE", "CANCELLED", "FAILED"}[s]
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
