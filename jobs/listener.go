package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
	"listening.to/orm"
	"listening.to/types"
	"log"
	"os"
	"time"
)

func main() {
	o, err := orm.New("sqlite3")
	if err != nil {
		log.Fatal("Failed to create orm: ", err)
	}
	defer o.Destroy()

	auth := spotify.NewAuthenticator(os.Getenv("REDIRECT_URI"), spotify.ScopeUserReadCurrentlyPlaying)
	rows, err := o.Query(&types.Account{})
	accs := rows.GetAccounts()
	if err != nil {
		log.Print(err)
	}

	for {
		for _, acc := range accs {
			time.Sleep(100 * time.Millisecond)
			c := auth.NewClient(acc.Token)
			cur, err := c.PlayerCurrentlyPlaying()
			if err != nil {
				switch err.Error() {
				case "API rate limit exceeded":
					log.Print(err)
					time.Sleep(3 * time.Second)
					continue
				default:
					log.Print(err)
					continue
				}
			}
			if cur.Item == nil {
				// podcasters lul
				continue
			}
			p := &types.Playing{
				acc.ID,
				cur.Item.ID,
				cur.Timestamp,
			}
			o.Write(p)
			// Make sure we keep track if a token changes
			t, err := c.Token()
			if err != nil {
				log.Print(err)
			}
			if acc.Token.AccessToken != t.AccessToken {
				if t.RefreshToken != "" {
					acc.Token = t
				} else {
					//Conserve current RefreshToken so that we don't get shut out
					acc.Token.AccessToken, acc.Token.Expiry = t.AccessToken, t.Expiry
				}
				go func(acc *types.Account) {
					o.Write(acc)
				}(acc)
			}
		}
	}
}
