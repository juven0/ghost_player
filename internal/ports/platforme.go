package ports

import "context"

type Platfome struct {
	Name      string
	StreamUrl string
}

type PlatformeInterface interface {
	Search(ctx context.Context, query string) ([]Track, error)
	StreamUrlFormat(streamUlr, id string) (string)
}


