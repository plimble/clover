package oauth2

type CreateAccessTokenRequest struct {
	ClientID  string
	UserID    string
	Scopes    []string
	ExpiresIn int
	Extras    map[string]interface{}
}

type TokenGenerator interface {
	CreateAccessToken(req *CreateAccessTokenRequest) (string, error)
	CreateRefreshToken() string
	CreateCode() string
}
