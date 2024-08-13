package dto

import "order-svc/pkg"

type AuthResponse struct {
	User  pkg.UserInfo `json:"user"`
	Token pkg.Token    `json:"token"`
}
