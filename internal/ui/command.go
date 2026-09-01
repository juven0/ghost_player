package ui

import (
	"player/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

// listenCmd relaie le flux d'événements player vers des messages TUI.
func listenCmd(player *domain.PlayerService) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-player.EventChan()
		if !ok {
			return nil
		}
		return playerEventMsg{event: ev}
	}
}