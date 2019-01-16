package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"listening.to/orm"
	"listening.to/types"
	"log"
	"net/http"
	"os"
	_ "time"
)

var o *orm.Orm

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	auth := spotify.NewAuthenticator(os.Getenv("REDIRECT_URI"), spotify.ScopeUserReadPrivate)
	auth.AuthURL("1")
	token, err := auth.Token("1", r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	client := auth.NewClient(token)
	user, err := client.CurrentUser()
	if err != nil {
		log.Print(err)
	}
	fmt.Fprintf(w, "%+v\n", user)

	a := types.Account{user.ID, token.AccessToken, token.RefreshToken}
	err = o.Write(a)
	if err != nil {
		log.Print(err)
	}
	log.Print("Wrote %v", a)

}

func listeningHandler(w http.ResponseWriter, r *http.Request) {
	auth := spotify.NewAuthenticator(os.Getenv("REDIRECT_URI"), spotify.ScopeUserReadCurrentlyPlaying)
	rows, err := o.Query(types.Account{})
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	var accs []types.Account
	for rows.Next() {
		var id, atok, rtok string
		rows.Scan(&id, &atok, &rtok)
		accs = append(accs, types.Account{id, atok, rtok})
	}

	for {
		for _, acc := range accs {
			c := auth.NewClient(&oauth2.Token{
				AccessToken:  acc.AccessToken,
				RefreshToken: acc.RefreshToken,
			})
			cur, err := c.PlayerCurrentlyPlaying()
			if err != nil {
				log.Print(err)
			}

			fmt.Fprintf(w, "%+v", cur.Item)
		}
	}

}

func main() {
	var err error
	o, err = orm.New("sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer o.Destroy()
	log.Printf("%+v\n", o)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/listening", listeningHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
