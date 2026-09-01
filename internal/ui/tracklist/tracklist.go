package tracklist

import (
	"context"
	"fmt"

	"player/internal/domain"
	"player/internal/ports"
	"player/internal/ui/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var searchKey = key.NewBinding(
	key.WithKeys("S"),
	key.WithHelp("S", "search in plateforme"),
)

// trackItem adapte ports.Track en list.Item.
type trackItem struct {
	track ports.Track
}

func (t trackItem) Title() string       { return t.track.Title }
func (t trackItem) Description() string { return t.track.Uploader }
func (t trackItem) FilterValue() string { return t.track.Title }

type searchCompleteMsg struct {
	tracks []ports.Track
	err    error
}

type playStartedMsg struct {
	track ports.Track
}

type playErrorMsg struct {
	err error
}

type playStoppedMsg struct{}

func searchCmd(ctx context.Context, track *domain.TrackService, query string, platformName string) tea.Cmd {
	return func() tea.Msg {
		tracks, err := track.Search(ctx, query, platformName)
		return searchCompleteMsg{tracks: tracks, err: err}
	}
}

func playCmd(ctx context.Context, player *domain.PlayerService, track ports.Track) tea.Cmd {
	return func() tea.Msg {
		if err := player.Play(ctx, track); err != nil {
			return playErrorMsg{err: err}
		}
		return playStartedMsg{track: track}
	}
}

type Model struct {
	list      list.Model
	input     textinput.Model
	track     *domain.TrackService
	player    *domain.PlayerService
	ctx       context.Context
	width     int
	height    int
	msg       string
	isSearch  bool
	isPlaying bool
	currentID string
	platform  string
}

func New(ctx context.Context, track *domain.TrackService, player *domain.PlayerService) Model {
	delegate := list.NewDefaultDelegate()
	ti := textinput.New()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Songs"
	return Model{
		list:   l,
		input:  ti,
		track:  track,
		player: player,
		ctx:    ctx,
	}
}

func (m Model) Init() tea.Cmd {
	return searchCmd(m.ctx, m.track, "shenseea", "")
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case searchCompleteMsg:
		if msg.err != nil {
			m.msg = fmt.Sprintf("❌ Erreur de recherche: %v", msg.err)
			return m, nil
		}
		items := make([]list.Item, len(msg.tracks))
		for i, t := range msg.tracks {
			items[i] = trackItem{track: t}
		}
		m.list.SetItems(items)
		m.msg = fmt.Sprintf("%d résultats trouvés", len(items))
		return m, nil

	case playStartedMsg:
		m.isPlaying = true
		m.currentID = msg.track.ID
		m.msg = "▶️  Lecture: " + msg.track.Title

	case playErrorMsg:
		m.isPlaying = false
		m.msg = "❌ Erreur de lecture: " + msg.err.Error()

	case playStoppedMsg:
		m.isPlaying = false
		m.currentID = ""
		m.msg = "⏹️  Lecture arrêtée"

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		if m.isSearch {
			switch msg.Type {
			case tea.KeyEnter:
				q := m.input.Value()
				if q != "" {
					m.isSearch = false
					m.msg = "🔍 Recherche en cours..."
					m.input.SetValue("")
					return m, searchCmd(m.ctx, m.track, q, m.platform)
				}
			case tea.KeyEsc:
				m.isSearch = false
				m.input.SetValue("")
				return m, nil
			}
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
		switch msg.Type {
		case tea.KeyEnter:
			if item, ok := m.list.SelectedItem().(trackItem); ok {
				return m, playCmd(m.ctx, m.player, item.track)
			}
		}
		if key.Matches(msg, searchKey) {
			m.isSearch = true
			m.input.Focus()
			return m, textinput.Blink
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	listHeight := height - 5
	if listHeight < 10 {
		listHeight = 10
	}
	m.list.SetSize(width, listHeight)
}

func (m Model) View() string {
	var view string

	if m.isSearch {
		view = "🔍 Rechercher sur YouTube:\n\n"
		view += m.input.View()
		view += "\n\n(Enter pour rechercher, Esc pour annuler)"
		if m.height > 0 {
			view = styles.AppStyle.Height(m.height).MaxHeight(m.height).Render(view)
		}
		return view
	}

	view = m.list.View()
	if m.msg != "" {
		view += "\n" + styles.AccentTextStyle.Render(m.msg)
	}
	return styles.AppStyle.Render(view)
}

func (m *Model) SetActivePlateform(name string) {
	m.platform = name
}
