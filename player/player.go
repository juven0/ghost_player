package player

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

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

func PlayCmd(url string) tea.Cmd {
	return func() tea.Msg {
		err := PlayAudio(url)
		if err != nil {
			return nil
		}
		return nil
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

	var videos []videoInfo
	scanner := bufio.NewScanner(strings.NewReader(result.Stdout))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

		var video videoInfo
		if err := json.Unmarshal([]byte(line), &video); err != nil {
			continue
		}

		videos = append(videos, video)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading output: %w", err)
	}

	if len(videos) == 0 {
		return nil, fmt.Errorf("no videos found")
	}

	return videos, nil
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

func getStreamURL(url string) (string, error) {
	cmd := exec.Command("yt-dlp", "-g", "-f", "bestaudio", url)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func PlayAudio(url string) error {
	streamURL, err := getStreamURL(url)
	if err != nil {
		return err
	}
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", streamURL)
	return cmd.Run()
}
