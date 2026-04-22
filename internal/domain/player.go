package domain

import (
	"context"
	"fmt"
	"player/internal/ports"
)

type PlayerService struct {
	player   ports.Player
	resolver ports.StreamURLResolver
	event    chan ports.PlayerEvent
}

func NewPlayerService(player ports.Player, resolver ports.StreamURLResolver) *PlayerService {
	return &PlayerService{
		player:   player,
		resolver: resolver,
		event:    make(chan ports.PlayerEvent),
	}
}

func (s *PlayerService) Play(context context.Context, track ports.Track) error {
	streamURL, err := s.resolver.Resolve(context, track)
	if err != nil {
		return fmt.Errorf("error to resolve stream url %s, %w", track.Title, err)
	}

	if err := s.player.StartPlay(context, streamURL); err != nil {
		return fmt.Errorf("error to play track %s, %w", track.Title, err)
	}

	return nil

}

func (s *PlayerService) Pause() error {
	return s.player.Pause()
}

func (s *PlayerService) Resume() error {
	return s.player.Resume()
}

func (s *PlayerService) Stop() error {
	return s.player.Stop()
}

func (s *PlayerService) SetVolume(volume int) error {
	return s.player.SetVolume(volume)
}

func (s *PlayerService) GetVolume() (int, error) {
	return s.player.GetVolume()
}

func (s *PlayerService) Seek(seconds float64, mode ports.SeekMode) error {
	return s.player.Seek(seconds, mode)
}

func (s *PlayerService) SetMute(mute bool) error {
	return s.player.SetMute(mute)
}

func (s *PlayerService) GetMute() (bool, error) {
	return s.player.GetMute()
}

func (s *PlayerService) Close() error {
	return s.player.Close()
}
