package tui

import (
	"player/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	footer      footer
	sidbare     plateformModel
	trackList   trackItemModel
	width       int
	height      int
	renderCount int
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
	listTitleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)
	mutedListTitleStyle = listTitleStyle.
				Background(mutedColor)
)

var (
	sidebarWidth = 25
	footerHeight = 2
)

func NewModel() Model {
	m := Model{
		footer:    newFooter(),
		sidbare:   newPlateformeList(),
		trackList: newTrackList(),
	}
	m.width = 80
	m.height = 24
	m.updateSizes()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.renderCount++
	var cmdList tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyLeft, tea.KeyRight:
			m.togglePanel()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
	}
	m.sidbare, cmdList = m.sidbare.Update(msg)
	return m, tea.Batch(cmdList)
}

func (m *Model) updateSizes() {
	// contentWidth := m.width - sidebarWidth - 4
	contentHeight := m.height - footerHeight - 12
	m.footer.SetSize(m.width-2, footerHeight)
	m.sidbare.SetSize(sidebarWidth, contentHeight)
}

func (m Model) View() string {
	bodyHeight := m.height - footerHeight

	body := lipgloss.JoinHorizontal(lipgloss.Left, styles.TrackBoxStyle.
		Width(m.width-2).
		Height(bodyHeight-50).
		Render(m.sidbare.View()), m.trackList.View())

	footer := m.footer.View()

	return lipgloss.JoinVertical(lipgloss.Left, body, footer)
}

func (m *Model) togglePanel() {
	m.sidbare.Focus()
}
