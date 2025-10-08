package player

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lrstanley/go-ytdlp"
)

type Player struct {
	ctx context.Context
}

type videoInfo struct {
	Title    string `json:"title"`
	Duration string `json:"duration"`
	URL      string `json:"url"`
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
