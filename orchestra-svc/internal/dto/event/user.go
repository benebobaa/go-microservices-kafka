package event

type UserRequest struct {
	Username string `json:"username"`
}

type UserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type User struct {
	Id            string `json:"id"`
	AccountBankId string `json:"account_bank_id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
}
