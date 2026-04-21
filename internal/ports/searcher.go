package ports

import (
	"context"
	"player/internal/domain"
)

type Searcher interface {
	Search(ctx context.Context, query string) ([]domain.Track, error)
}

type StreamURLResolver interface {
	Resolve(ctx context.Context, videoID string) (string, error)
}
