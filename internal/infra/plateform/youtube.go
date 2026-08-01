package plateform

import (
	"fmt"

	"player/internal/ports"
)

type Youtube struct {
	Platforme ports.Platforme
}

func NewYoutube() ports.PlatformeInterface {
	return &Youtube{
		Platforme: ports.Platforme{
			Name:      "Youtube",
			StreamUrl: "https://www.youtube.com/watch?v=",
			SearchURL: "ytsearch",
			Color:     "DEFAULT",
		},
	}
}

func (p *Youtube) StreamUrlFormat(id string) string {
	return fmt.Sprintf("%s%s", p.Platforme.SearchURL, id)
}

func (p *Youtube) FormatQuery(query string, max int) (string, error) {
	return fmt.Sprintf("%s%d:%s", p.Platforme.StreamUrl, max, query), nil
}
