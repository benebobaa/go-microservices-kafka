package dto

type Payment struct {
	RefId     string  `json:"ref_id"`
	Amount    float64 `json:"amount"`
	AccountId string  `json:"account_id"`
}

type RefundRequest struct {
	RefId string `json:"ref_id"`
}

type PaymentRequest struct {
	RefId     string  `json:"ref_id"`
	Amount    float64 `json:"amount"`
	AccountID string  `json:"account_id"`
}
