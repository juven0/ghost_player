package mpv

import (
	"context"
	"fmt"

	"player/internal/ports"
)

func (m *MpvPlayer) ObserveEvents(ctx context.Context) <-chan ports.PlayerEvent {
	ch := make(chan ports.PlayerEvent, 20)

	go func() {
		defer close(ch)

		properties := []string{"percent-pos", "pause", "volume", "mute", "duratrion", "time-pose"}

		for i, prop := range properties {
			m.client.SendCommand([]string{"observer_property", fmt.Sprintf("%d", i+1), prop})
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:

			}
		}
	}()

	return ch
}
