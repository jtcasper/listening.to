package main

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"listening.to/orm"
	"listening.to/types"
	"log"
	"os"
)

func main() {
	o, err := orm.New("sqlite3")
	if err != nil {
		log.Fatal("Failed to create orm: ", err)
	}

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatal("Failed to exchange token: ", err)
	}

	client := spotify.Authenticator{}.NewClient(token)

	rows, err := o.RawQuery(
		`SELECT *
    FROM PLAYING
    WHERE NOT EXISTS (
      SELECT NULL
      FROM TRACK
      WHERE TRACK.ID == PLAYING.TRACK_ID
    )`,
	)
	if err != nil {
		log.Fatal(err)
	}

	plays := rows.GetPlaying()
	tracks := plays.UniqueTracks()
	var fTracks []*spotify.FullTrack

	for len(tracks.Tracks) > 0 {
		fts, err := client.GetTracks(func(tc *types.TrackContainer) []spotify.ID {
			var ids []spotify.ID
			for len(ids) <= 50 && len(tracks.Tracks) > 0 {
				var t *types.Track
				t, tracks.Tracks = tracks.Tracks[0], tracks.Tracks[1:]
				ids = append(ids, t.ID)
			}
			return ids
		}(tracks)...)
		if err != nil {
			log.Print("Failed to retrieve tracks: ", err)
		}

		fTracks = append(fTracks, fts...)
	}
	for _, ft := range fTracks {
		err = o.Write(&types.Track{ft.ID, ft.Album.ID, ft.Name, ft.Duration})
	}
}
