package ui

import (
	"player/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type panel int

const (
	panelSibare panel = iota
	panelTracklist
	panelFooter
)

var (
	sidebarWidth = 25
	footerHeight = 2
)

func (m *UIModel) View() string {
	bodyHeight := m.height - footerHeight - 4
	if bodyHeight < 15 {
		bodyHeight = 15
	}

	body := styles.TrackBoxStyle.
		Width(m.width - 2).
		Height(bodyHeight).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, m.sidebar.View(), m.tracklist.View()))

	return lipgloss.JoinVertical(lipgloss.Left, body, m.footer.View())
}

func (m *UIModel) toggelPannel(k tea.KeyType) {
	if k == tea.KeyRight {
		m.active = (m.active + 1) % 3
	}
	if k == tea.KeyLeft {
		if m.active > 0 {
			m.active--
		} else {
			m.active = 2
		}
	}
}

func (m *UIModel) updateSize() {
	contentWidth := m.width - sidebarWidth - 4
	contentHeight := m.height - footerHeight - 6
	bodyHeight := m.height - footerHeight - 4

	if bodyHeight < 15 {
		bodyHeight = 15
	}

	m.footer.SetSize(m.width-2, footerHeight)
	m.sidebar.SetSize(sidebarWidth, contentHeight)
	m.tracklist.SetSize(contentWidth, bodyHeight)
}
