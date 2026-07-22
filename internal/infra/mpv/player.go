package mpv

import (
	"context"

	"player/internal/ports"
)

type MpvPlayer struct {
	client  *Client
	process *Process
}

func NewMpvPlayer(client *Client, process *Process) ports.Player {
	return &MpvPlayer{
		client:  client,
		process: process,
	}
}

func (m *MpvPlayer) Connect(ctx context.Context) error {
	return m.client.pipe.Connect(ctx, m.process.PipePath())
}
