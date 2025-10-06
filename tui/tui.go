package tui

import (
	"player/styles"

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
	footerHeight = 5
)

func NewModel() Model {
	m := Model{
		footer:  NewFooter(),
		sidbare: newPlateformeList(),
	}
	// Initialiser avec une taille par défaut
	m.width = 80
	m.height = 24
	m.updateSizes()
	return m
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
	bodyHeight := m.height - footerHeight

	body := styles.TrackBoxStyle.
		Width(m.width).
		Height(bodyHeight - 10).
		Render(m.sidbare.View())

	footer := m.footer.View()

	return lipgloss.JoinVertical(lipgloss.Left, body, footer)
}
