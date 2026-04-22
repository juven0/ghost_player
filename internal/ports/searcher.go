package ports

import (
	"context"
)

type Searcher interface {
	Search(ctx context.Context, query string) ([]Track, error)
}

type StreamURLResolver interface {
	Resolve(ctx context.Context, track Track) (string, error)
}
