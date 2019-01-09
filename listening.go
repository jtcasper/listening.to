package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	_ "time"
)

type AuthResponseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int16  `json:"expires_in"`
	Scope        string `json:"scope"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if errorCode, exists := values["error"]; exists {
		fmt.Fprintf(w, "There was an error signing up: %s", errorCode)
		log.Printf("Error: %s", errorCode)
	}

	resp, err := http.PostForm("https://accounts.spotify.com/api/token",
		url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {values.Get("code")},
			"redirect_uri":  {os.Getenv("REDIRECT_URI")},
			"client_id":     {os.Getenv("CLIENT_ID")},
			"client_secret": {os.Getenv("CLIENT_SECRET")},
		},
	)
	defer resp.Body.Close()

	var auth = AuthResponseBody{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	json.Unmarshal(body, &auth)

	if err != nil {
		log.Print(err)
	}

	userResp, err := doAPIRequest(auth.AccessToken, "me")
	if err != nil {
		log.Print(err)
	}
	defer userResp.Body.Close()

	var user spotify.User
	userBody, err := ioutil.ReadAll(userResp.Body)
	if err != nil {
		log.Print(err)
	}

	json.Unmarshal(userBody, &user)
	fmt.Fprintf(w, "%+v\n", user)

	db, err := sql.Open("sqlite3", "listening.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	res, err := db.Exec("INSERT OR REPLACE INTO ACCOUNTS (ID, ACCESS_TOKEN, REFRESH_TOKEN) VALUES ($1, $2, $3) ",
		user.ID,
		auth.AccessToken,
		auth.RefreshToken,
	)
	if err != nil {
		log.Print(err)
	}
	fmt.Fprintf(w, "%+v", res)

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
	db, err := sql.Open("sqlite3", "listening.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM ACCOUNTS")
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	var accs []Account
	for rows.Next() {
		var id, atok, rtok string
		rows.Scan(&id, &atok, &rtok)
		accs = append(accs, Account{id, atok, rtok})
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

			json.Unmarshal(b, &cur)
			if err != nil {
				log.Print(err)
			}
			fmt.Fprintf(w, "%+v", cur.Item)
		}
	}

}

type Account struct {
	ID           string
	AccessToken  string
	RefreshToken string
}

func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/listening", listeningHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
