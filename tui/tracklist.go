package tui

import (
	"player/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type plateformItem struct {
	name string
}

func (p plateformItem) FilterValue() string {
	return p.name
}

var plateforms = []plateformItem{
	{name: "Youtube"},
	{name: "Spotifye"},
	{name: "Deezer"},
}

type plateformModel struct {
	list    list.Model
	width   int
	height  int
	focused bool
}

func newPlateformeList() plateformModel {
	l := list.New(plateformsToListItem(plateforms), list.NewDefaultDelegate(), 0, 0)
	l.Title = "Plateforme"
	l.DisableQuitKeybindings()
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	return plateformModel{
		list: l,
	}
}

func (m plateformModel) Init() tea.Cmd {
	return func() tea.Msg {
		return m.list.Items()[0]
	}
}

func (m plateformModel) Update(msg tea.Msg) (plateformModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			item := m.list.SelectedItem()
			m.list.FilterInput.SetValue("")
			return m, func() tea.Msg { return item }
		}
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m plateformModel) View() string {
	style := styles.FocusedStyle
	if !m.focused {
		style = mutedPanelStyle
	}
	return style.
		Width(m.width).
		Height(m.height).
		Render(m.list.View())
}

func (m *plateformModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

func plateformsToListItem(sitems []plateformItem) []list.Item {
	items := make([]list.Item, len(sitems))
	for i, item := range sitems {
		items[i] = list.Item(item)
	}
	return items
}
