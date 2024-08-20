package dto

type UserValidateRequest struct {
	Username string `json:"username"`
}

type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateBankIDRequest struct {
	Username      string `json:"username"`
	AccountBankID string `json:"account_bank_id"`
}

type UserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	AccountBankID string `json:"account_bank_id"`
	Email         string `json:"email"`
}
