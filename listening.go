package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
	_ "golang.org/x/oauth2"
	"listening.to/orm"
	"listening.to/types"
	"log"
	"net/http"
	"os"
)

var o *orm.Orm

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "content/index.html")
	} else {
		http.ServeFile(w, r, r.URL.Path[1:])
	}
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	st := r.FormValue("state")
	if st == "" {
		http.NotFound(w, r)
		log.Print("Couldn't get token")
		return
	}
	auth := spotify.NewAuthenticator(os.Getenv("REDIRECT_URI"), spotify.ScopeUserReadPrivate)
	token, err := auth.Token(st, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	client := auth.NewClient(token)
	user, err := client.CurrentUser()
	if err != nil {
		log.Print(err)
	}

	a := types.Account{spotify.ID(user.ID), token}

	err = o.Write(a)
	if err != nil {
		log.Print(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "account_info",
		Value: string(a.ID),
		Path:  "/",
		//		Domain:   "2600:1700:24d1:4b50::4d8",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: false,
	})
	fmt.Fprintf(w, "%+v\n", user)

}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	accountCookie, err := r.Cookie("account_info")
	if err != nil {
		log.Print(err)
		return
	}

	trackRows, err := o.RawQuery(
		`select track.*, count(id)
		from track
		join playing on track.id = playing.track_id
		where playing.account_id = $1
		group by id;`,
		accountCookie.Value,
	)
	if err != nil {
		log.Print(err)
		return
	}

	trackContainer := trackRows.GetTracksWithCounts()
	b, err := json.Marshal(trackContainer)
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	var err error
	o, err = orm.New("sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer o.Destroy()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/analyze", analyzeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
