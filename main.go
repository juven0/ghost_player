package main

import (
	"context"
	"log"

	"player/internal/infra/ipc"
	"player/internal/infra/mpv"
	"player/tui"

	tea "github.com/charmbracelet/bubbletea"
	//	"github.com/lrstanley/go-ytdlp"
)

func main() {
	// ytdlp.MustInstall(context.TODO(), nil)

	ctx := context.Background()
	transport := ipc.NewPipeWindows()
	client := mpv.NewClient(transport)
	process := mpv.NewProcess(ctx)
	_ = mpv.NewMpvPlayer(client, process)

	m := tui.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
