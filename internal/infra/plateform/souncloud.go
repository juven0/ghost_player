package plateform

import (
	"fmt"

	"player/internal/ports"
)

type SounCloud struct {
	Platforme ports.Platforme
}

func NewSounCloud() ItemPlateforme {
	return ItemPlateforme{
		Platforme: ports.Platforme{
			Name:      "SounCloud",
			StreamUrl: "https://www.souncloud.com/",
			SearchURL: "scsearch10:",
			Color:     "DEFAULT",
		},
		Platformer: &SounCloud{},
	}
}

func (p *SounCloud) StreamUrlFormat(track ports.Track) string {
	return track.SourceURL
}

func (p *SounCloud) FormatQuery(query string, max int) (string, error) {
	return fmt.Sprintf("scsearch10:%s", query), nil
}

func (p *SounCloud) TrackMaping(data map[string]interface{}) (ports.Track, error) {
	return ports.Track{
		ID:         getString(data, "id"),
		Title:      getString(data, "title"),
		Source:     ports.TrackSourceSoundCloud,
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
