package dto

type BankRegistrationRequest struct {
	Username string  `json:"username" valo:"notblank"`
	Email    string  `json:"email" valo:"notblank"`
	Deposit  float64 `json:"deposit" valo:"min=1000"`
}

type BankRegistrationUpdate struct {
	CustomerID string `json:"customer_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Status     string `json:"-"`
}
