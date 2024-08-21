package interfaces

import "context"

type HtppRequest interface {
	POST(ctx context.Context, path string, req interface{}, res interface{}) error
}
