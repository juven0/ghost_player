package ui

import (
	"player/internal/ports"
)

// playerEventMsg est le message interne du fan-in du flux d'événements player.
type playerEventMsg struct {
	event ports.PlayerEvent
}