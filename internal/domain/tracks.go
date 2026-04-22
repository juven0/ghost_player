package domain

import (
	"context"
	"errors"
	"fmt"
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
	tracks, err := s.searcher.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error searching tracks: %w", err)
	}
	return tracks, nil
}
