package types

import (
	"github.com/zmb3/spotify"
)

type (
	// Wraps necessary spotify.Track fields for our ORM
	Track struct {
		ID       spotify.ID `json:"id"`
		AlbumID  spotify.ID `json:"album_id,omitempty"`
		Name     string     `json:"name,omitempty"`
		Duration int        `json:"duration,omitempty"`
	}
	TrackContainer struct {
		Tracks []*Track `json:"track_container"`
	}
	TrackCount struct {
		Track
		PlayCount int `json:"play_count,omitempty"`
	}
	TrackCountContainer struct {
		TrackCounts []*TrackCount `json:"track_container"`
	}
)

func (t *Track) Table() string {
	return "TRACK"
}
