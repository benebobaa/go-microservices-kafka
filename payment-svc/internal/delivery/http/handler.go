package http

import (
	"payment-svc/internal/usecase"
)

type Handler struct {
	u *usecase.Usecase
}

func NewHandler(u *usecase.Usecase) *Handler {
	return &Handler{
		u: u,
	}
}
