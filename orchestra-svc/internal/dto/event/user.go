package event

type UserRequest struct {
	Username string `json:"username"`
}

type UserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
