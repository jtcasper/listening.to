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

	a := types.Account{spotify.ID(user.ID), token}

	err = o.Write(a)
	if err != nil {
		log.Print(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "account_info",
		Value:    string(a.ID),
		Path:     "/",
		Domain:   "2600:1700:24d1:4b50::6f3",
		Secure:   false,
		HttpOnly: false,
	})
	fmt.Fprintf(w, "%+v\n", user)

}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("account_info")
	if err != nil {
		log.Print(err)
		return
	}
	rows, err := o.RawQuery("SELECT * FROM PLAYING WHERE ACCOUNT_ID = $1", cookie.Value)
	if err != nil {
		log.Print(err)
		return
	}
	pc := rows.GetPlaying()
	log.Print(pc)
	log.Print(pc.MostPlayed())

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
