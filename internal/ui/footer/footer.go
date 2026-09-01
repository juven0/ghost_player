package footer

import (
	"player/internal/ports"
	"player/internal/ui/styles"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	spinner       spinner.Model
	progress      progress.Model
	width         int
	height        int
	progressValue float64
	playing       bool
}

func New() *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return &Model{
		spinner:  s,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.spinner, cmd = m.spinner.Update(msg)
	return *m, cmd
}

func (m *Model) SetEvent(ev ports.PlayerEvent) {
	m.handleEvent(ev)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	progressWidth := width - 13
	if progressWidth > 0 {
		m.progress.Width = progressWidth
	}
}

func (m *Model) handleEvent(ev ports.PlayerEvent) {
	switch ev.Type {
	case ports.EventFileLoaded:
		m.playing = true
	case ports.EventPausedChanged:
		m.playing = !ev.Pauased
	case ports.EventFileEnded:
		m.playing = false
		m.progressValue = 1
	}
	// progression (percent) propagée par EventFileLoaded/TimePos
}

func (m Model) View() string {
	playButton := styles.ActiveButtonStyle.Padding(0, 1).Margin(0).Render(styles.IconPlay)

	style := styles.MutedPanelStyle.
		Padding(0, 1).
		Width(m.width).
		Height(m.height)
	return style.
		Render(
			playButton,
			styles.TrackProgressStyle.Width(m.width).Render(
				m.progress.ViewAs(m.progressValue),
			),
		)
}
