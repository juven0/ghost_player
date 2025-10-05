package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type footer struct {
	spinner  spinner.Model
	progress progress.Model
	width    int
	height   int
}

func NewFooter() footer {
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
	style := panelStyle.
		Width(m.width).
		Height(m.height).
		Padding(0, 1)

	content := mutedTextStyle.Render(" Ananas ") + "\n" + m.progress.ViewAs(0.25)

	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		style.Render(content),
	)
}

func (m *footer) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.progress.Width = w - 4
}
