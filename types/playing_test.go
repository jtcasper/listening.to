package types_test

import (
	"github.com/zmb3/spotify"
	"listening.to/types"
	"testing"
)

func TestUniqueTracksReturnsUnique(t *testing.T) {
	pc := types.PlayingContainer{
		[]*types.Playing{
			&types.Playing{TrackID: spotify.ID(1)},
			&types.Playing{TrackID: spotify.ID(1)},
			&types.Playing{TrackID: spotify.ID(2)},
		},
	}
	tc := pc.UniqueTracks()
	if len(tc.Tracks) != 2 {
		t.Fail()
	}
}
