package response

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
