package ports

import "context"

type Player interface {
	Connect(ctx context.Context) error
	GetPercentPos() (int, error)
	Play() error
	Pause() error
	Stop() error
	Resume() error
	SetVolume(volume int) error
	GetVolume() (int, error)
	Seek(seconds float64, mode SeekMode) error
	SetMute(mute bool) error
	GetMute() (bool, error)
	StartPlay(ctx context.Context, streamURL string) error
	Close() error
}

type SeekMode string

const (
	SeekModeAbsolute SeekMode = "absolute"
	SeekModeRelative SeekMode = "relative"
)

type PlayerEvent struct {
	Type    EventType
	TimePos float64
	Volume  int
	Mute    bool
	Pauased bool
	Reason  EndReason
	Error   error
}

type EventType int

const (
	EventFileLoaded EventType = iota
	EventFileEnded
	EventVolumeChanged
	EventMuteChanged
	EventPausedChanged
	EventError
	EventShutdown
	EventIdle
)

type EndReason string

const (
	EndReasonEof     EndReason = "eof"
	EndReasonError   EndReason = "error"
	EndReasonQuit    EndReason = "quit"
	EndReasonStop    EndReason = "stop"
	EndReasonUnknown EndReason = "unknown"
)
