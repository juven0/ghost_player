package ui

import (
	"context"

	"player/internal/domain"
)

type Deps struct {
	player *domain.PlayerService
	track  *domain.TrackService
	cancel context.CancelFunc
}

func NewDeps(player *domain.PlayerService, track *domain.TrackService, cancel context.CancelFunc) Deps {
	return Deps{
		player: player,
		track:  track,
		cancel: cancel,
	}
}