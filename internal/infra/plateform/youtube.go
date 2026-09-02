package plateform

import (
	"fmt"

	"player/internal/ports"
)

type Youtube struct {
	Platforme ports.Platforme
}

func NewYoutube() ItemPlateforme {
	return ItemPlateforme{
		Platforme: ports.Platforme{
			Name:      "Youtube",
			StreamUrl: "https://www.youtube.com/watch?v=",
			SearchURL: "ytsearch",
			Color:     "DEFAULT",
		},
		Platformer: &Youtube{},
	}
}

func (p *Youtube) StreamUrlFormat(id string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
}

func (p *Youtube) FormatQuery(query string, max int) (string, error) {
	return fmt.Sprintf("ytsearch%d:%s", max, query), nil
}
