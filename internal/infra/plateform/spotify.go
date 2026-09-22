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

func (p *Spotify) StreamUrlFormat(track ports.Track) string {
	return fmt.Sprintf("https://open.spotify.com/track/%s", track.ID)
}

func (p *Spotify) FormatQuery(query string, max int) (string, error) {
	return fmt.Sprintf("ytsearch%d:%s", max, query), nil
}

func (p *Spotify) TrackMaping(data map[string]interface{}) (ports.Track, error) {
	return ports.Track{
		ID:         getString(data, "id"),
		Title:      getString(data, "title"),
		Source:     ports.TrackSourceSpotify,
		Duration:   getFloat64(data, "duration"),
		Uploader:   getString(data, "uploader"),
		SourceURL:  getString(data, "webpage_url"),
		Artist:     getString(data, "artist"),
		Album:      getString(data, "album"),
		ViewCount:  getInt64(data, "view_count"),
		UploadDate: getString(data, "upload_date"),
		SteamURL:   "",
	}, nil
}
