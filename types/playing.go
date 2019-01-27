package types

import (
	"github.com/zmb3/spotify"
)

// Wraps spotify.CurrentlyPlaying for use by our ORM.
type Playing struct {
	CP        spotify.CurrentlyPlaying
	AccountID spotify.ID
}

func (p *Playing) Table() string {
	return "PLAYING"
}
