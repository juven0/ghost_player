package memory

import (
	"sync"

	"player/internal/ports"
)

type Tracks struct {
	mu        sync.Mutex
	liked     map[string]ports.Track
	playlists map[string][]ports.Track
}

func NewTracks() ports.Tracks {
	return &Tracks{
		liked:     make(map[string]ports.Track),
		playlists: make(map[string][]ports.Track),
	}
}

func (t *Tracks) Like(track *ports.Track) ([]ports.Track, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.liked[track.ID] = *track

	liked := make([]ports.Track, 0, len(t.liked))
	for _, tr := range t.liked {
		liked = append(liked, tr)
	}
	return liked, nil
}

func (t *Tracks) PPlaylist(playlistID string) ([]ports.Track, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if tracks, ok := t.playlists[playlistID]; ok {
		return tracks, nil
	}
	return nil, nil
}

func (t *Tracks) NewPlaylist(name string, tracks []ports.Track) (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.playlists[name] = tracks
	return name, nil
}

func (t *Tracks) DeletePlaylist(playlistID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.playlists, playlistID)
	return nil
}