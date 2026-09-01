package sidebar

import (
	"player/internal/ports"
	"player/internal/ui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type platformItem struct {
	name string
}

type PlatformSelectedMsg struct {
	Name string
}

func (p platformItem) FilterValue() string { return p.name }

type Model struct {
	list    list.Model
	width   int
	height  int
	focused bool
}

func New(platforms []ports.Platforme) Model {
	items := make([]platformItem, len(platforms))
	for i, p := range platforms {
		items[i] = platformItem{name: p.Name}
	}
	l := list.New(platformsToItems(items), newPlatformDelegate(false), 0, 0)
	l.Title = "Plateforme"
	l.DisableQuitKeybindings()
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	return Model{list: l}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		if len(m.list.Items()) == 0 {
			return nil
		}
		item := m.list.Items()[0].(platformItem)
		return PlatformSelectedMsg{Name: item.name}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
		item, ok := m.list.SelectedItem().(platformItem)
		if !ok {
			return m, nil
		}
		return m, func() tea.Msg {
			return PlatformSelectedMsg{Name: item.name}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	style := styles.FocusedStyle
	if !m.focused {
		style = styles.MutedPanelStyle
	}
	return style.Width(m.width).Height(m.height).Render(m.list.View())
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height-6)
}

func platformsToItems(in []platformItem) []list.Item {
	out := make([]list.Item, len(in))
	for i, p := range in {
		out[i] = list.Item(p)
	}
	return out
}
