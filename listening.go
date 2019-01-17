package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
	_ "golang.org/x/oauth2"
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

	a := types.Account{user.ID, token}

	err = o.Write(a)
	if err != nil {
		log.Print(err)
	}

}

func listeningHandler(w http.ResponseWriter, r *http.Request) {
	auth := spotify.NewAuthenticator(os.Getenv("REDIRECT_URI"), spotify.ScopeUserReadCurrentlyPlaying)
	rows, err := o.Query(types.Account{})
	accs := rows.GetAccounts()
	if err != nil {
		log.Print(err)
	}

	for {
		for index, acc := range accs {
			c := auth.NewClient(acc.Token)
			cur, err := c.PlayerCurrentlyPlaying()
			if err != nil {
				switch err.Error() {
				case "The access token expired":
					t, err := c.Token()
					if err != nil {
						log.Print(err)
					}
					acc.Token = t
					accs[index] = acc
					o.Write(acc)
					continue
				default:
					continue
				}
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
