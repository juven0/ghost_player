package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	footer  footer
	sidbare plateformModel
	width   int
	height  int
}

var (
	mutedColor = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	mutedPanelStyle = panelStyle.
			BorderForeground(mutedColor)

	mutedTextStyle = lipgloss.NewStyle().
			Foreground(mutedColor)
)

var (
	sidebarWidth = 25
	footerHeight = 2
)

func NewModel() Model {
	return Model{
		footer:  NewFooter(),
		sidbare: newPlateformeList(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
	}

	return m, nil
}

func (m Model) updateSizes() {
	m.footer.SetSize(m.width, footerHeight)
}

func (m Model) View() string {
	body := m.sidbare.View()
	footer := m.footer.View()

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Top,
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Top,
			body,
			lipgloss.Place(
				m.width,
				footerHeight,
				lipgloss.Bottom,
				lipgloss.Left,
				footer,
			),
		),
	)
}
