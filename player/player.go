package player

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/go-ytdlp"
)

type Player struct {
	cancel context.CancelFunc
	cmd    *exec.Cmd
}

type playProgress struct {
	CurrentTime string
	TotalTime   string
	percentage  int
}
type PlayStoppedMsg struct{}

var currentPlayer *Player

type PlayStartedMsg struct {
	VideoID string
	Title   string
}

type PlayErrorMsg struct {
	Err error
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

func NewPlayer() *Player {
	_, cancel := context.WithCancel(context.Background())
	return &Player{
		cancel: cancel,
	}
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

func PlayCmd(video videoInfo) tea.Cmd {
	return func() tea.Msg {
		if currentPlayer != nil {
		}

		streamURL, err := getStreamURL(video.ID)
		if err != nil {
			return PlayErrorMsg{Err: err}
		}

		ctx, cancel := context.WithCancel(context.Background())

		cmd := exec.CommandContext(ctx, "mpv",
			"--no-video",
			"--really-quiet",
			"--keep-open=no",
			streamURL,
		)
		currentPlayer = &Player{
			cmd:    cmd,
			cancel: cancel,
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil
		}

		progressRgx := regexp.MustCompile(`A:\s+\d{2}:(\d{2}:\d{2})\s+/\s+\d{2}:(\d{2}:\d{2})\s+\((\d+)%\)`)

		go func() {
			if err := cmd.Run(); err != nil && ctx.Err() == nil {
				fmt.Printf("Play error: %v\n", err)
			}
			currentPlayer = nil
		}()

		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				matches := progressRgx.FindStringSubmatch(line)
				if len(matches) > 3 {
					// progress := strconv.Itoa(matches[3])
				}
			}
		}()

		return PlayStartedMsg{
			VideoID: video.ID,
			Title:   video.Title,
		}
	}
}

func StopAudio() {
	if currentPlayer != nil {
		if currentPlayer.cancel != nil {
			currentPlayer.cancel()
		}
		if currentPlayer.cmd != nil && currentPlayer.cmd.Process != nil {
			currentPlayer.cmd.Process.Kill()
		}
		currentPlayer = nil
	}
}

func StopCmd() tea.Cmd {
	return func() tea.Msg {
		StopAudio()
		return PlayStoppedMsg{}
	}
}

func IsPlaying() bool {
	return currentPlayer != nil
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

func getStreamURL(mediaId string) (string, error) {
	ctx := context.Background()
	mediaURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", mediaId)

	result, err := ytdlp.New().
		Format("bestaudio").
		GetURL().
		NoWarnings().
		Run(ctx, mediaURL)
	if err != nil {
		return "", fmt.Errorf("failed to get stream URL: %w", err)
	}

	streamURL := strings.TrimSpace(result.Stdout)
	if streamURL == "" {
		return "", fmt.Errorf("empty stream URL")
	}

	return streamURL, nil
}

func PlayAudio(mediaId string) error {
	streamURL, err := getStreamURL(mediaId)
	if err != nil {
		return err
	}
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", streamURL)
	return cmd.Run()
}
