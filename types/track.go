package types

import (
	"github.com/zmb3/spotify"
)

// Wraps spotify.Track for our ORM
type Track struct {
	Track *spotify.FullTrack
}

func (t *Track) Table() string {
	return "TRACK"
}
