package types

import (
	"github.com/zmb3/spotify"
)

type (
	// Wraps necessary spotify.Track fields for our ORM
	Track struct {
		ID       spotify.ID
		AlbumID  spotify.ID
		Name     string
		Duration int
	}
	TrackContainer struct {
		Tracks []*Track
	}
)

func (t *Track) Table() string {
	return "TRACK"
}
