package interfaces

import "context"

type Client interface {
	GET(ctx context.Context, url string, request interface{}, response interface{}) error
}

type Producer interface {
	SendMessage(key string, message []byte) error
}
