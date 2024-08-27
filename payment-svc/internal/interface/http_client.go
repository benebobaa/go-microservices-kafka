package interfaces

import "context"

type Client interface {
	GET(ctx context.Context, url string, request interface{}, response interface{}) error
	POST(ctx context.Context, url string, request interface{}, response interface{}) error
	PATCH(ctx context.Context, url string, request interface{}, response interface{}) error
}
