package ports

type Platforme struct {
	Name      string
	StreamUrl string
	SearchURL string
	Color     string
}

type PlatformeInterface interface {
	FormatQuery(query string, max int) (string, error)
	StreamUrlFormat(id string) string
}
