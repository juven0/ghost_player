package player

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/go-ytdlp"
)
const (
	Loading = iota
	Stopped = iota
	Paused  = iota
	Playing = iota
)

type Player struct {
	cancel context.CancelFunc
	cmd    *exec.Cmd
	ctx    context.Context
	info   PlayerInfo
	ch     chan PlayerMsg
	pipe    string
	state   int
	stream  string
}

type PlayerInfo struct {
	Duration string
	Current  string
	Progress int
}
type PlayStoppedMsg struct{}
type PlayerErrorMsg error
type PlayerStateChangedMsg string
type PlayerOutputMsg string

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
	Info videoInfo
}

func (t TrackItem) Title() string       { return t.Info.Title }
func (t TrackItem) Description() string { return t.Info.Uploader }
func (t TrackItem) FilterValue() string { return t.Info.Title }

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
		ch:     make(chan PlayerMsg, 10),
		state:  Stopped,
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

func (p *Player) PlayCmd(video videoInfo) {
		if currentPlayer != nil {
		}

		streamURL, err := getStreamURL(video.ID)
		if err != nil {
		}

		p.pipe = path.Join(os.TempDir(), "mpvsocket")

		p.ctx, p.cancel = context.WithCancel(context.Background())

		p.cmd = exec.CommandContext(p.ctx, "mpv",
			streamURL,
			"--no-video",
			"--ytdl-format=bestaudio",
			fmt.Sprintf("--input-ipc-server=%s", p.pipe),
			"--quiet",
		)

			stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		p.ch <- PlayerErrorMsg(fmt.Errorf("error creating stdout pipe: %v", err))
		return
	}

	progressRegex := regexp.MustCompile(`A:\s+\d{2}:(\d{2}:\d{2})\s+/\s+\d{2}:(\d{2}:\d{2})\s+\((\d+)%\)`)

	if err := p.cmd.Start(); err != nil {
		p.ch <- PlayerErrorMsg(fmt.Errorf("Error starting command: %v\n", err))
		return
	}

	p.setState(Loading)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			matches := progressRegex.FindStringSubmatch(line)
			if len(matches) > 3 {
				if p.state == Loading {
					p.setState(Playing)
				}
				progress, _ := strconv.Atoi(matches[3])
				if progress != p.info.Progress && progress < 100 {
					p.info = PlayerInfo{
						Current:  matches[1],
						Duration: matches[2],
						Progress: progress,
					}
					select {
					case p.ch <- PlayerProgressMsg(p.info):
					default:
					}
				}
			} else {
				p.ch <- PlayerOutputMsg(line)
			}
		}
	}()

	go func() {
		defer os.Remove(p.pipe)
		err := p.cmd.Wait()
		if err != nil {
			p.setState(Stopped)
			return
		}
		p.info = PlayerInfo{
			Progress: 100,
			Current:  p.info.Duration,
			Duration: p.info.Duration,
		}
		p.ch <- PlayerProgressMsg(p.info)
	}()
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
			Info: video,
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

func (p *Player) setState(state int) {
	if p.state == state {
		return
	}
	p.state = state
	p.ch <- PlayerStateChangedMsg(state)
}

func (p *Player) sendSocket(command string) error {
	conn, err := net.Dial("unix", p.pipe)
	if err != nil {
		return fmt.Errorf("Error connecting to socket: %v\n", err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		return fmt.Errorf("Error writing to socket: %v\n", err)
	}
	reader := bufio.NewReader(conn)
	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Failed to read response: %v\n", err)
	}
	return nil
}

func (p *Player) VideoToListItem(videos []videoInfo) []list.Item {
	items := make([]list.Item, len(videos))
	for i, v := range videos {
		items[i] = TrackItem{Info: v}
	}
	return items
}