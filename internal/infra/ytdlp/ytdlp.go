package ytdlp

import (
	"context"
	"fmt"
	"strings"

	"github.com/lrstanley/go-ytdlp"
)

type Ytdlp struct{}

func (yt Ytdlp) Resolve(ctx context.Context, url string) (string, error) {
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
