package tui

import (
	"fmt"
	"io"

	"player/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PlatfomDeleget struct {
	normalStyle   lipgloss.Style
	selectedStyle lipgloss.Style
}

func (d PlatfomDeleget) Height() int                               { return 1 }
func (d PlatfomDeleget) Spacing() int                              { return 0 }
func (d PlatfomDeleget) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func newSimpleListDelegate(focused bool) PlatfomDeleget {
	d := newDefaultListDelegate(focused)
	return PlatfomDeleget{
		normalStyle:   d.Styles.NormalTitle,
		selectedStyle: d.Styles.SelectedTitle,
	}
}

func newDefaultListDelegate(focused bool) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	if focused {
		return d
	}
	d.Styles.SelectedTitle = d.Styles.NormalTitle
	d.Styles.SelectedDesc = d.Styles.NormalTitle
	d.Styles.NormalTitle = d.Styles.NormalDesc
	return d
}

func (d PlatfomDeleget) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var content string
	if m.Index() == index {
		content = styles.TrackListActiveStyle.Render(item.FilterValue())
	} else {
		content = styles.TrackListStyle.Render(item.FilterValue())
	}

	fmt.Fprint(w, content)
}
