package types

type (
	AuthResponseBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int16  `json:"expires_in"`
		Scope        string `json:"scope"`
	}
	APIError struct {
		ErrorContainer struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		} `json:"error"`
	}
)
