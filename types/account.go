package types

import (
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type Account struct {
	ID    spotify.ID
	Token *oauth2.Token
}

func (a *Account) String() string {
	return fmt.Sprintf("ID: %s\nAccessToken: %s\nRefreshToken: %s\nExpiry: %s",
		a.ID,
		a.Token.AccessToken,
		a.Token.RefreshToken,
		a.Token.Expiry,
	)
}

func (a *Account) Table() string {
	return "ACCOUNT"
}

func NewAccount() *Account {
	return &Account{Token: &oauth2.Token{}}
}
