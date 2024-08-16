package dto

type ProductQuantityRetryRequest struct {
	Quantity   int    `json:"quantity" valo:"min=1"`
	EventID    string `json:"event_id" valo:"notblank"`
	InstanceID string `json:"instance_id" valo:"notblank"`
}

type ProductReserveRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
