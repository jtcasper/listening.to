package types

import (
	"github.com/zmb3/spotify"
)

// Wraps spotify.CurrentlyPlaying for use by our ORM.
type (
	Playing struct {
		AccountID spotify.ID
		TrackID   spotify.ID
		Timestamp int64
	}
	PlayingContainer struct {
		Plays []*Playing
	}
)

func (p *Playing) Table() string {
	return "PLAYING"
}

func (pc *PlayingContainer) MostPlayed() *Playing {
	var maxCount int
	var maxPlay *Playing
	playCounts := make(map[spotify.ID]int)
	for _, p := range pc.Plays {
		if cnt, ok := playCounts[p.TrackID]; ok {
			cnt += 1
			playCounts[p.TrackID] = cnt
			if cnt > maxCount {
				maxPlay = p
				maxCount = cnt
			}
		} else {
			playCounts[p.TrackID] = 1
		}
	}
	return maxPlay
}
