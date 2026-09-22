package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	AccentColor       = lipgloss.Color("#ba299aff")
	BackgroundColor   = lipgloss.Color("#6b6b6b")
	ActiveTextColor   = lipgloss.Color("#EEE")
	NormalTextColor   = lipgloss.Color("#CCC")
	InactiveTextColor = lipgloss.Color("#888")
	BorderMutedColor  = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}
)

var (
	IconPlay     = "▶"
	IconStop     = "■"
	IconLiked    = "💛"
	IconNotLiked = "🤍"
)

var AccentTextStyle = lipgloss.NewStyle().Foreground(AccentColor)

// PanelStyle est le cadre uniforme de tous les panneaux (sidebar, tracklist,
// footer). Padding(0,1) + bordure arrondie => frame (x=4, y=2).
var (
	PanelStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderMutedColor)

	PanelFocusedStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(AccentColor)
)

// PanelFrameSize retourne la taille du cadre (horizontal, vertical) compris
// dans Width/Height, i.e. la part à soustraire pour obtenir le contenu.
func PanelFrameSize() (x, y int) {
	return PanelStyle.GetFrameSize()
}

var (
	TrackListStyle = lipgloss.NewStyle().
			Padding(0, 1)
	TrackListActiveStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(AccentColor)
	TrackTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#dcdcdc")).
			Bold(true)
)

var (
	TrackVersionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#999999"))
	TrackArtistStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#dcdcdc"))
	TrackAddInfoStyle = lipgloss.NewStyle().
				Align(lipgloss.Right).
				Width(26)
)

var (
	StatusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065"))
)