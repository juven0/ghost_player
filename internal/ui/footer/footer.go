package footer

import (
	"player/internal/ports"
	"player/internal/ui/styles"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// playButtonReserve est la largeur du bouton play (3) plus son séparateur (1).
const playButtonReserve = 4

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

	frameX, _ := styles.PanelFrameSize()
	progressWidth := width - frameX - playButtonReserve
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
	playButtonStyle := lipgloss.NewStyle().
		Background(styles.AccentColor).
		Foreground(styles.ActiveTextColor).
		Padding(0, 1)
	playButton := playButtonStyle.Render(styles.IconPlay)

	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		playButton,
		" ",
		m.progress.ViewAs(m.progressValue),
	)

	return styles.PanelStyle.Width(m.width - 2).Height(m.height).Render(content)
}
