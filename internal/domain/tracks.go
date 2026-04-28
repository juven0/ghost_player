package domain

import (
	"context"
	"errors"
	"fmt"
	"player/internal/ports"
)

type TrackService struct {
	searcher ports.Searcher
	resolver ports.StreamURLResolver
	Traks    ports.Tracks
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

func (s *TrackService) ResolveStreamURL(ctx context.Context, track ports.Track, resolver ports.StreamURLResolver) (string, error) {
	return resolver.Resolve(ctx, track)
}

func (s *TrackService) Like(track *ports.Track) ([]ports.Track, error) {
	return s.Traks.Like(track)
}

func (s *TrackService) PPlaylist(playlistID string) ([]ports.Track, error) {
	return s.Traks.PPlaylist(playlistID)
}

func (s *TrackService) NewPlaylist(name string, tracks []ports.Track) (string, error) {
	return s.Traks.NewPlaylist(name, tracks)
}

func (s *TrackService) DeletePlaylist(playlistID string) error {
	return s.Traks.DeletePlaylist(playlistID)
}
