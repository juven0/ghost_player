package tui

import (
	"player/styles"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type endMsg struct{}

type footer struct {
	spinner  spinner.Model
	progress progress.Model
	width    int
	height   int
}

func newFooter() footer {
	return footer{
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m footer) Init() tea.Cmd {
	return nil
}

func (m footer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m footer) View() string {
	playButton := styles.ActiveButtonStyle.Padding(0, 1).Margin(0).Render(styles.IconPlay)

	// progressBar := m.progress.View()
	// content := styles.TrackProgressStyle.Width(m.width).Render(progressBar)
	// content = lipgloss.JoinHorizontal(lipgloss.Top, playButton, content)

	style := panelStyle.
		Padding(0, 1).
		Width(m.width).
		Height(m.height)
	return style.
		Render(
			playButton,
			styles.TrackProgressStyle.Width(m.width).Render(m.progress.View()),
		)

	// return styles.TrackBoxStyle.Width(m.width).Render(content)
}

func (m *footer) SetSize(w, h int) {
	m.width = w
	m.height = h

	progressWidth := w - 13

	if progressWidth > 0 {
		m.progress.Width = progressWidth
	}
}

func endCmd() tea.Msg {
	return endMsg{}
}