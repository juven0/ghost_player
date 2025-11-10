package player

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/go-ytdlp"
)

type Player struct {
	cancel context.CancelFunc
	cmd    *exec.Cmd
	ctx    context.Context
	info   PlayerInfo
	ch     chan PlayerMsg
}

type PlayerInfo struct {
	Current  string
	Total    string
	Progress int
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
type (
	PlayerMsg         interface{}
	PlayerProgressMsg PlayerInfo
)

func NewPlayer() *Player {
	ctx, cancel := context.WithCancel(context.Background())
	return &Player{
		cancel: cancel,
		ctx:    ctx,
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

func (p *Player) PlayCmd(video videoInfo) tea.Cmd {
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
			"--term-status-msg=A:%{=time-pos}/%{=duration} (%{percent-pos}%)",
			streamURL,
		)

		//cmd := exec.Command("mpv",
		//"--no-video",
		//"--script-opts=ytdl_hook-ytdl_path=yt-dlp", // 👈 indique à mpv d'utiliser yt-dlp
		//"--ytdl-format=bestaudio/best",
		//streamURL,
		//)
		currentPlayer = &Player{
			cmd:    cmd,
			cancel: cancel,
			ctx:    ctx,
			ch:     make(chan PlayerMsg, 10),
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return PlayErrorMsg{Err: err}
		}

		// progressRgx := regexp.MustCompile(`A:\s+\d{2}:(\d{2}:\d{2})\s+/\s+\d{2}:(\d{2}:\d{2})\s+\((\d+)%\)`)

		go func() {
			if err := cmd.Run(); err != nil && ctx.Err() == nil {
				fmt.Printf("Play error: %v\n", err)
			}
			currentPlayer = nil
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			fmt.Println(scanner)
			for scanner.Scan() {
				line := scanner.Text()
				matches := regexp.MustCompile(`A:(\d+\.\d+)/(\d+\.\d+) \((\d+)%\)`).FindStringSubmatch(line)
				if len(matches) == 4 {
					current := matches[1]
					total := matches[2]
					progress, _ := strconv.Atoi(matches[3])
					fmt.Println(progress)
					info := PlayerInfo{
						Current:  current,
						Total:    total,
						Progress: progress,
					}

					select {
					case currentPlayer.ch <- PlayerProgressMsg(info):
					default:
					}
				}
			}
		}()

		return PlayStartedMsg{
			VideoID: video.ID,
			Title:   video.Title,
		}
	}
}

func (p *Player) Ch() chan PlayerMsg {
	return p.ch
}

func (p *Player) Info() PlayerInfo {
	return p.info
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
		Format("bestaudio/best").
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
