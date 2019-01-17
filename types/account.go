package types

import (
	"golang.org/x/oauth2"
)

type Account struct {
	ID    string
	Token *oauth2.Token
}

func (a Account) Table() string {
	return "ACCOUNT"
}
