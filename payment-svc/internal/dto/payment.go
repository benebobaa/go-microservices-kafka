package dto

type Transaction struct {
	Id            string  `json:"id"`
	RefId         string  `json:"ref_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	AccountBankId string  `json:"account_bank_id"`
}
type RefundRequest struct {
	RefId string `json:"ref_id"`
}

type PaymentRequest struct {
	RefId         string  `json:"ref_id"`
	Amount        float64 `json:"amount"`
	AccountBankID string  `json:"account_bank_id"`
}

type AccountBalanceRequest struct {
	Username string  `json:"username"`
	Deposit  float64 `json:"deposit"`
}

type AccountBalance struct {
	AccountBankID string  `json:"account_bank_id"`
	Balance       float64 `json:"balance"`
	Username      string  `json:"username"`
}
