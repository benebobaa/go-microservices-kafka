package dto

type BaseResponse[T any] struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error,omitempty"`
	Data       T      `json:"data,omitempty"`
}
