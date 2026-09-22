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

const (
	sidebarWidth = 25
	footerHeight = 1
)

func (m *UIModel) View() string {
	body := lipgloss.JoinHorizontal(lipgloss.Top, m.sidebar.View(), m.tracklist.View())
	return lipgloss.JoinVertical(lipgloss.Top, body, m.footer.View())
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
	_, frameY := styles.PanelFrameSize()

	bodyHeight := m.height - footerHeight
	if bodyHeight < 4+frameY {
		bodyHeight = 4 + frameY
	}

	m.footer.SetSize(m.width, footerHeight)
	m.sidebar.SetSize(sidebarWidth, bodyHeight)
	m.tracklist.SetSize(m.width-sidebarWidth, bodyHeight)
}
