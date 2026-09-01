package plateform

import "player/internal/ports"

var youtube = NewYoutube()

var Gostplatforme = []ports.Platforme{
	youtube.Platforme,
}
