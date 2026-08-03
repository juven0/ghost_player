package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type deps struct{}

type panel int

const (
	panelSibare panel = iota
	panelTracklist
	panelFooter
)

type UIModel struct {
	deps        deps
	width       int
	height      int
	renderCount int
	ctx         context.Context
	active      panel
}

func NewModel(deps deps) UIModel {
	m := UIModel{
		deps: deps,
	}
	m.active = panelTracklist
	m.width = 80
	m.height = 24
	return m
}

func (m UIModel) Init() tea.Cmd {
	return tea.Batch()
}

func (m *UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyLeft, tea.KeyRight:
			m.toggelPannel(msg)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *UIModel) View() string {
	return ""
}

func (m *UIModel) toggelPannel(t tea.Msg) {
	if t == tea.KeyRight {
		m.active = (m.active + 1) % 3
	}
	if t == tea.KeyLeft {
		if m.active > 0 {
			m.active--
		} else {
			m.active = 2
		}
	}
}

func (m *UIModel) updateSize() {
}
