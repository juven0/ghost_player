package ytdlp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"player/internal/ports"

	"github.com/lrstanley/go-ytdlp"
)

type Ytdlp struct{}

func NewYtdlp(ctx context.Context) ports.StreamResolver {
	return &Ytdlp{}
}

func (yt *Ytdlp) Resolve(ctx context.Context, url string) (string, error) {
	result, err := ytdlp.New().
		Format("bestaudio/best").
		GetURL().
		NoWarnings().
		Run(ctx, url)
	if err != nil {
		return "", fmt.Errorf("failed to get stream URL: %w", err)
	}

	streamURL := strings.TrimSpace(result.Stdout)
	if streamURL == "" {
		return "", fmt.Errorf("empty stream URL")
	}
	return streamURL, nil
}

func (yt *Ytdlp) Search(ctx context.Context, searchQuery string, maxRes int) ([]ports.Track, error) {
	dl := ytdlp.New().FlatPlaylist().DumpJSON()
	res, err := dl.Run(ctx, searchQuery)
	if err != nil {
		return []ports.Track{}, err
	}

	var trakRes []ports.Track
	scanner := bufio.NewScanner(strings.NewReader(res.Stdout))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

		var trak ports.Track
		if err := json.Unmarshal([]byte(line), &trak); err != nil {
			continue
		}

		trakRes = append(trakRes, trak)
	}

	if err := scanner.Err(); err != nil {
		return trakRes, fmt.Errorf("error reading output: %w", err)
	}

	if len(trakRes) == 0 {
		return trakRes, fmt.Errorf("no item found")
	}

	return trakRes, nil
}
