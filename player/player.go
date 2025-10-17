package player

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/go-ytdlp"
)

type Player struct {
	ctx context.Context
}

type videoInfo struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
	Uploader string  `json:"uploader"`
	URL      string  `json:"url"`
}

type TrackItem struct {
	Video videoInfo
}

func (t TrackItem) Title() string       { return t.Video.Title }
func (t TrackItem) Description() string { return t.Video.Uploader }
func (t TrackItem) FilterValue() string { return t.Video.Title }

type SearchCompleteMsg struct {
	Results []videoInfo
	Err     error
}

func SearchYTCmd(query string, maxRes int) tea.Cmd {
	return func() tea.Msg {
		results, err := SearchYoutube(query, maxRes)
		return SearchCompleteMsg{
			results,
			err,
		}
	}
}

func SearchYoutube(query string, maxResult int) ([]videoInfo, error) {
	ctx := context.Background()

	dl := ytdlp.New().FlatPlaylist().DumpJSON()

	searchQuery := fmt.Sprintf("ytsearch%d:%s", maxResult, query)
	result, err := dl.Run(ctx, searchQuery)
	if err != nil {
		return []videoInfo{}, fmt.Errorf("search failed: %w", err)
	}

	var data struct {
		Entries []videoInfo `json:"entries"`
	}

	err = json.Unmarshal([]byte(result.Stdout), &data)
	if err != nil {
		return []videoInfo{}, fmt.Errorf("error parse results: %w", err)
	}

	return data.Entries, nil
}

func VideoToListeItem(videos []videoInfo) []list.Item {
	items := make([]list.Item, len(videos))
	for i, video := range videos {
		items[i] = TrackItem{
			Video: video,
		}
	}
	return items
}
