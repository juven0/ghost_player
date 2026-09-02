package domain

import (
	"context"
	"errors"
	"fmt"
	"player/internal/infra/plateform"
	"player/internal/ports"
)

type TrackService struct {
	Platforms []plateform.ItemPlateforme
	resolver  ports.StreamResolver
	Traks     ports.Tracks
}

func NewTrack(platform []plateform.ItemPlateforme, resolver ports.StreamResolver, tracks ports.Tracks) *TrackService {
	return &TrackService{
		Platforms: platform,
		resolver:  resolver,
		Traks:     tracks,
	}
}

func (s *TrackService) Search(ctx context.Context, query string, platformName string) ([]ports.Track, error) {
	if query == "" {
		return nil, errors.New("empty query")
	}

	platform, err := getPlateformByName(platformName, s.Platforms)
	if err != nil {
		return nil, fmt.Errorf("error finding platform: %w", err)
	}

	formattedQuery, err := platform.FormatQuery(query, 10)
	if err != nil {
		return nil, fmt.Errorf("error formatting query: %w", err)
	}

	tracks, err := s.resolver.Search(ctx, formattedQuery, 10)
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

func getPlateformByName(name string, platforms []plateform.ItemPlateforme) (*plateform.ItemPlateforme, error) {
	if name == "" {
		return &platforms[0], nil
	}
	for _, p := range platforms {
		if p.Platforme.Name == name {
			return &p, nil
		}
	}
	return nil, errors.New("platform not found")
}
