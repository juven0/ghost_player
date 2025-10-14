package tui

import (
	"player/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type trackItemModel struct {
	list         list.Model
	width        int
	height       int
	msg          string
	keys         *trackKeyMap
	delegateKeys *delegateKeyMap
}

type trackKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

func newListeKeyMap() *trackKeyMap {
	return &trackKeyMap{
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

	const numItem = 10
	traks := make([]list.Item, numItem)

	delegate := newTrackDelegate(delegateKey)
	tracks := list.New(traks, delegate, 0, 0)
	tracks.Title = "Songs"
	tracks.Styles.Title = styles.TitleStyle
	tracks.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			trakKey.toggleSpinner,
			trakKey.toggleStatusBar,
			trakKey.toggleTitleBar,
			trakKey.toggleHelpMenu,
			trakKey.togglePagination,
		}
	}

	return trackItemModel{
		list:         tracks,
		keys:         trakKey,
		delegateKeys: delegateKey,
	}
}
