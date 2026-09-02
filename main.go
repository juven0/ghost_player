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
	"player/internal/ports"
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
	tracks := memory.NewTracks()

	// youtube := plateform.NewYoutube()
	spotify := plateform.NewSpotify()

	var Gostplatforme = []plateform.ItemPlateforme{
		// youtube,
		spotify,
	}

	platforms := make([]ports.Platforme, len(Gostplatforme))
	for i, p := range Gostplatforme {
		platforms[i] = p.Platforme
	}

	playerService := domain.NewPlayerService(player, resolver, Gostplatforme)
	trackService := domain.NewTrack(Gostplatforme, resolver, tracks)

	deps := ui.NewDeps(playerService, trackService, nil, platforms)
	m := ui.NewModel(deps)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
