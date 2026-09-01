package ui

import (
	"context"

	"player/internal/ui/footer"
	"player/internal/ui/sidebar"
	"player/internal/ui/tracklist"

	tea "github.com/charmbracelet/bubbletea"
)

type UIModel struct {
	deps        Deps
	width       int
	height      int
	renderCount int
	active      panel

	sidebar   sidebar.Model
	tracklist tracklist.Model
	footer    footer.Model

	ctx context.Context
}

func NewModel(deps Deps) *UIModel {
	ctx, cancel := context.WithCancel(context.Background())
	deps.cancel = cancel

	m := &UIModel{
		deps:   deps,
		ctx:    ctx,
		active: panelTracklist,
		width:  80,
		height: 24,
	}
	m.sidebar = sidebar.New()
	m.tracklist = tracklist.New(ctx, deps.track, deps.player)
	m.footer = *footer.New()
	return m
}

func (m *UIModel) Init() tea.Cmd {
	return tea.Batch(
		m.tracklist.Init(),
		m.footer.Init(),
		m.sidebar.Init(),
		listenCmd(m.deps.player),
	)
}

func (m *UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.renderCount++
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyLeft, tea.KeyRight:
			m.toggelPannel(msg.Type)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSize()

	case playerEventMsg:
		m.footer.SetEvent(msg.event)
		cmds = append(cmds, listenCmd(m.deps.player))
	}

	var cmd tea.Cmd
	m.sidebar, cmd = m.sidebar.Update(msg)
	cmds = append(cmds, cmd)
	m.tracklist, cmd = m.tracklist.Update(msg)
	cmds = append(cmds, cmd)
	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
