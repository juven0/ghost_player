package ui

import (
	"context"

	"player/internal/domain"
	"player/internal/ports"
)

type Deps struct {
	player     *domain.PlayerService
	track      *domain.TrackService
	cancel     context.CancelFunc
	plateforms []ports.Platforme
}

func NewDeps(player *domain.PlayerService, track *domain.TrackService, cancel context.CancelFunc, plateforms []ports.Platforme) Deps {
	return Deps{
		player:     player,
		track:      track,
		cancel:     cancel,
		plateforms: plateforms,
	}
}
