package domain

import (
	"context"
	"errors"
	"fmt"
	"player/internal/ports"
)

type TrackService struct {
	Platform ports.PlatformeInterface
	resolver ports.StreamResolver
	Traks    ports.Tracks
}

func NewTrack(platform ports.PlatformeInterface, resolver ports.StreamResolver, tracks ports.Tracks) *TrackService {
	return &TrackService{
		Platform: platform,
		resolver: resolver,
		Traks:    tracks,
	}
}

func (s *TrackService) Search(ctx context.Context, query string) ([]ports.Track, error) {
	if query == "" {
		return nil, errors.New("empty query")
	}
	tracks, err := s.resolver.Search(ctx, query, 10)
	if err != nil {
		return nil, fmt.Errorf("error searching tracks: %w", err)
	}
	return tracks, nil
}

func (s *TrackService) ResolveStreamURL(ctx context.Context, url string, resolver ports.StreamResolver) (string, error) {
	return s.resolver.Resolve(ctx, url)
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
