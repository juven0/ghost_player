package plateform

import (
	"fmt"

	"player/internal/ports"
)

type Spotify struct {
	Platforme ports.Platforme
}

func NewSpotify() ItemPlateforme {
	return ItemPlateforme{
		Platforme: ports.Platforme{
			Name:      "Spotify",
			StreamUrl: "https://open.spotify.com/track/",
			SearchURL: "spotify",

			Color: "DEFAULT",
		},
		Platformer: &Spotify{},
	}
}

func (p *Spotify) StreamUrlFormat(id string) string {
	return fmt.Sprintf("https://open.spotify.com/track/%s", id)
}

func (p *Spotify) FormatQuery(query string, max int) (string, error) {
	return fmt.Sprintf("ytsearch%d:%s", max, query), nil
}
