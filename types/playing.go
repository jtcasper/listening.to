package types

import (
	"github.com/zmb3/spotify"
)

// Wraps spotify.CurrentlyPlaying for use by our ORM.
type Playing struct {
	AccountID spotify.ID
	TrackID   spotify.ID
	Timestamp int64
}

func (p *Playing) Table() string {
	return "PLAYING"
}
