package plateform

import "player/internal/ports"

type ItemPlateforme struct {
	Platforme  ports.Platforme
	Platformer ports.PlatformeInterface
}

func (i ItemPlateforme) StreamUrlFormat(id string) string {
	return i.Platformer.StreamUrlFormat(id)
}

func (i ItemPlateforme) FormatQuery(query string, max int) (string, error) {
	return i.Platformer.FormatQuery(query, max)
}
