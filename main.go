package main

import (
	"context"
	"log"
	"player/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	
	ytdlp.MustInstall(context.TODO(), nil)
	m := tui.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
