package types

type Account struct {
	ID           string
	AccessToken  string
	RefreshToken string
}

func (a Account) Table() string {
	return "ACCOUNTS"
}
