package dto

type AuthRequest struct {
	Username string `json:"username" valo:"notblank,sizeMin=6"`
}
