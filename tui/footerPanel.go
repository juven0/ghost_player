package tui

import (
	"player/player"
	"player/styles"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type endMsg struct{}

type footer struct {
	player   *player.Player
	spinner  spinner.Model
	progress progress.Model
	width    int
	height   int
}

func newFooter(p *player.Player) footer {
	return footer{
		player:   p,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m footer) Init() tea.Cmd {
	return tea.Batch(
		m.listenCmd,
	)
}

func (m footer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case player.PlayerProgressMsg:
		if msg.Progress == 100 {
			return m, tea.Batch(
				m.listenCmd,
				endCmd,
			)
		}
		return m, m.listenCmd
	case player.PlayerMsg:
		return m, m.listenCmd
	}
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
			styles.TrackProgressStyle.Width(m.width).Render(m.progress.ViewAs(float64(m.player.Info().Progress)/100)),
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

func (m *footer) listenCmd() tea.Msg {
	return <-m.player.Ch()
}
