package domain

type Track struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Source     TrackSource `json:"source"`
	Duration   float64     `json:"duration"`
	Uploader   string      `json:"uploader"`
	SourceURL  string      `json:"url"`
	Artist     string      `json:"artist"`
	Album      string      `json:"album"`
	ViewCount  int64       `json:"view_count"`
	UploadDate string      `json:"upload_date"`
	SteamURL   string      `json:"stream_url"`
}

type TrackSource string

const (
	TrackSourceYouTube    TrackSource = "youtube"
	TrackSourceSoundCloud TrackSource = "soundcloud"
	TrackSourceSpotify    TrackSource = "spotify"
	TrackSourceLocal      TrackSource = "local"
)
