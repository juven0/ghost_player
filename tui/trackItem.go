package tui

import (
	"fmt"

	"player/player"
	"player/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type trackItemModel struct {
	list         list.Model
	input        textinput.Model
	width        int
	height       int
	msg          string
	keys         *trackKeyMap
	delegateKeys *delegateKeyMap
	isSearch     bool
}

type trackKeyMap struct {
	search           key.Binding
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

func newListeKeyMap() *trackKeyMap {
	return &trackKeyMap{
		search: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "search in plateforme"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

func newTrackList() trackItemModel {
	var (
		delegateKey = newDelegateKeyMap()
		trakKey     = newListeKeyMap()
	)

	ti := textinput.New()
	const numItem = 2
	traks := make([]list.Item, numItem)

	delegate := newTrackDelegate(delegateKey)
	tracks := list.New(traks, delegate, 0, 0)
	tracks.Title = "Songs"
	tracks.Styles.Title = styles.TitleStyle
	tracks.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			trakKey.search,
			trakKey.toggleSpinner,
			trakKey.toggleStatusBar,
			trakKey.toggleTitleBar,
			trakKey.toggleHelpMenu,
			trakKey.togglePagination,
		}
	}

	return trackItemModel{
		list:         tracks,
		input:        ti,
		keys:         trakKey,
		delegateKeys: delegateKey,
		isSearch:     false,
	}
}

func (m trackItemModel) Init() tea.Cmd {
	return player.SearchYTCmd("shenseea", 3)
}

func (m trackItemModel) Update(msg tea.Msg) (trackItemModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.AppStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		if m.isSearch {
			switch msg.Type {
			case tea.KeyEnter:
				query := m.input.Value()
				if query != "" {
					m.isSearch = false
					m.msg = "🔍 Recherche en cours..."
					m.input.SetValue("")
					return m, player.SearchYTCmd(query, 10)
				}
			case tea.KeyEsc:
				m.isSearch = false
				m.input.SetValue("")
				return m, nil
			}
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
		switch {
		case key.Matches(msg, m.keys.search):
			m.isSearch = true
			m.input.Focus()
			return m, textinput.Blink
		}
	case player.SearchCompleteMsg:
		if msg.Err != nil {
			m.msg = fmt.Sprintf("Erreur: %v", msg.Err)
			return m, nil
		}

		items := player.VideoToListeItem(msg.Results)
		m.list.SetItems(items)
		m.msg = fmt.Sprintf("%d résultats trouvés", len(items))

		return m, nil
	}
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel

	return m, cmd
}

func (m *trackItemModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, 5)
}

func (m trackItemModel) View() string {
	var view string

	if m.isSearch {
		view = "🔍 Rechercher sur YouTube:\n\n"
		view += m.input.View()
		view += "\n\n(Enter pour rechercher, Esc pour annuler)"
		return styles.AppStyle.Render(view)
	}

	view = m.list.View()

	if m.msg != "" {
		view += "\n" + styles.AccentTextStyle.Render(m.msg)
	}

	return styles.AppStyle.Render(view)
}
