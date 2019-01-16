package main

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
	"io/ioutil"
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

func doAPIRequest(token, endpoint string) (r *http.Response, err error) {

	req, err := newRequestWithToken(token, endpoint)
	if err != nil {
		return
	}
	r, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	return

}

func newRequestWithToken(token, endpoint string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", "https://api.spotify.com/v1/"+endpoint, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return
}

func listeningHandler(w http.ResponseWriter, r *http.Request) {
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
			r, err := doAPIRequest(acc.AccessToken, "me/player/currently-playing")
			if err != nil {
				log.Print(err)
			}

			var cur spotify.CurrentlyPlaying
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Print(err)
			}
			stat, err := checkAPIResponse(b)
			if err != nil {
				log.Print(err)
			}
			switch stat {
			case 401:
				//auth, err := requestAccessToken(acc.RefreshToken)
				//if err != nil {
				//	log.Print(err)
				//}
				//acc.AccessToken, acc.RefreshToken = auth.AccessToken, auth.RefreshToken
				//o.Write(acc)
				//log.Print("Updated account!!!")
				//continue
			}

			json.Unmarshal(b, &cur)
			if err != nil {
				log.Print(err)
			}
			fmt.Fprintf(w, "%+v", cur.Item)
		}
	}

}

func checkAPIResponse(b []byte) (stat int, err error) {
	var a types.APIError
	json.Unmarshal(b, &a)
	if err != nil {
		return
	}
	return a.ErrorContainer.Status, nil

	return stat, errors.New("No Error response from API found")
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
