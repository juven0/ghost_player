package domain

import (
	"context"
	"errors"
	"player/internal/ports"
)

type TrackService struct {
	searcher ports.Searcher
}

func NewTrack(searcher ports.Searcher) *TrackService {
	return &TrackService{
		searcher: searcher,
	}
}

func (s *TrackService) Search(ctx context.Context, query string) ([]ports.Track, error) {
	if query == "" {
		return nil, errors.New("empty query")
	}
	return s.searcher.Search(ctx, query)
}
