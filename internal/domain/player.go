package domain

import (
	"context"
	"fmt"
	"player/internal/infra/plateform"
	"player/internal/ports"
)

type PlayerService struct {
	Player    ports.Player
	Resolver  ports.StreamResolver
	Event     chan ports.PlayerEvent
	Platforms []plateform.ItemPlateforme
}

func NewPlayerService(player ports.Player, resolver ports.StreamResolver, platforms []plateform.ItemPlateforme) *PlayerService {
	return &PlayerService{
		Player:    player,
		Resolver:  resolver,
		Event:     make(chan ports.PlayerEvent),
		Platforms: platforms,
	}
}

func (s *PlayerService) EventChan() <-chan ports.PlayerEvent {
	return s.Event
}

func (s *PlayerService) Play(context context.Context, track ports.Track, platformName string) error {
	platform, err := getPlateformByName(platformName, s.Platforms)
	if err != nil {
		return fmt.Errorf("error resolving platform: %w", err)
	}

	streamURL, err := s.Resolver.Resolve(context, platform.StreamUrlFormat(track.ID))
	if err != nil {
		return fmt.Errorf("error to resolve stream url %s, %w", track.Title, err)
	}

	if err := s.Player.StartPlay(context, streamURL); err != nil {
		return fmt.Errorf("error to play track %s, %w", track.Title, err)
	}

	return nil

}

func (s *PlayerService) Pause() error {
	return s.Player.Pause()
}

func (s *PlayerService) Resume() error {
	return s.Player.Resume()
}

func (s *PlayerService) Stop() error {
	return s.Player.Stop()
}

func (s *PlayerService) SetVolume(volume int) error {
	return s.Player.SetVolume(volume)
}

func (s *PlayerService) GetVolume() (int, error) {
	return s.Player.GetVolume()
}

func (s *PlayerService) Seek(seconds float64, mode ports.SeekMode) error {
	return s.Player.Seek(seconds, mode)
}

func (s *PlayerService) SetMute(mute bool) error {
	return s.Player.SetMute(mute)
}

func (s *PlayerService) GetMute() (bool, error) {
	return s.Player.GetMute()
}

func (s *PlayerService) Close() error {
	return s.Player.Close()
}
