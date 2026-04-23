package ports

import (
	"context"
)

type StreamURLResolver interface {
	Resolve(ctx context.Context, url string) (string, error)
}
