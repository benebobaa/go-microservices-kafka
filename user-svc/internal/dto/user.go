package dto

type UserValidateRequest struct {
	Username string `json:"username"`
}

type UserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	AccountBankID string `json:"account_bank_id"`
	Email         string `json:"email"`
}
