package ports

import (
	"context"
)

type StreamResolver interface {
	Resolve(ctx context.Context, url string) (string, error)
	Search(ctx context.Context, searchQuery string, maxRes int) ([]Track, error)
}
