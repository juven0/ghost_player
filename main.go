package main

import (
	"context"
	"log"

	"player/internal/domain"
	"player/internal/infra/ipc"
	"player/internal/infra/memory"
	"player/internal/infra/mpv"
	"player/internal/infra/plateform"
	"player/internal/infra/ytdlp"
	ui "player/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx := context.Background()

	transport := ipc.NewPipeWindows()
	client := mpv.NewClient(transport)
	process := mpv.NewProcess(ctx)
	player := mpv.NewMpvPlayer(client, process)

	resolver := ytdlp.NewYtdlp(ctx)
	platform := plateform.NewYoutube()
	tracks := memory.NewTracks()

	playerService := domain.NewPlayerService(player, resolver)
	trackService := domain.NewTrack(platform.Platforme, resolver, tracks)

	deps := ui.NewDeps(playerService, trackService, nil)
	m := ui.NewModel(deps)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
