package plateform

import "player/internal/ports"

type ItemPlateforme struct {
	Platforme  ports.Platforme
	Platformer ports.PlatformeInterface
}

func (i ItemPlateforme) StreamUrlFormat(track ports.Track) string {
	return i.Platformer.StreamUrlFormat(track)
}

func (i ItemPlateforme) FormatQuery(query string, max int) (string, error) {
	return i.Platformer.FormatQuery(query, max)
}

func getString(data map[string]any, key string) string {
	value, ok := data[key]
	if !ok || value == nil {
		return ""
	}

	text, ok := value.(string)
	if !ok {
		return ""
	}

	return text
}

func getFloat64(data map[string]any, key string) float64 {
	value, ok := data[key]
	if !ok || value == nil {
		return 0
	}

	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	}

	return 0
}

func getInt64(data map[string]any, key string) int64 {
	value, ok := data[key]
	if !ok || value == nil {
		return 0
	}

	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case float32:
		return int64(v)
	}

	return 0
}
