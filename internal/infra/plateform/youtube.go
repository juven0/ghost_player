package plateform

import (
	"fmt"

	"player/internal/ports"
)

func NewYoutube() ItemPlateforme {
	return ItemPlateforme{
		Platforme: ports.Platforme{
			Name:      "Youtube",
			StreamUrl: "https://www.youtube.com/watch?v=",
			SearchURL: "ytsearch",
			Color:     "DEFAULT",
		},
	}
}

func (p *ItemPlateforme) StreamUrlFormat(id string) string {
	return fmt.Sprintf("%s%s", p.Platforme.SearchURL, id)
}

func (p *ItemPlateforme) FormatQuery(query string, max int) (string, error) {
	return fmt.Sprintf("ytsearch%d:%s", max, query), nil
}
