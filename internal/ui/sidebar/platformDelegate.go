package sidebar

import (
	"fmt"
	"io"

	"player/internal/ui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type platformDelegate struct {
	normalStyle   lipgloss.Style
	selectedStyle lipgloss.Style
}

func newPlatformDelegate(focused bool) platformDelegate {
	d := newDefaultDelegate(focused)
	return platformDelegate{
		normalStyle:   d.Styles.NormalTitle,
		selectedStyle: d.Styles.SelectedTitle,
	}
}

func newDefaultDelegate(focused bool) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	if focused {
		return d
	}
	d.Styles.SelectedTitle = d.Styles.NormalTitle
	d.Styles.SelectedDesc = d.Styles.NormalTitle
	d.Styles.NormalTitle = d.Styles.NormalDesc
	return d
}

func (d platformDelegate) Height() int                               { return 1 }
func (d platformDelegate) Spacing() int                              { return 0 }
func (d platformDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d platformDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var content string
	if m.Index() == index {
		content = styles.TrackListActiveStyle.Render(item.FilterValue())
	} else {
		content = styles.TrackListStyle.Render(item.FilterValue())
	}
	fmt.Fprint(w, content)
}
