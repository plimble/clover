package oauth2

//go:generate mockery -name Storage
type Storage interface {
	GetClient(id string) (*Client, error)
	GetClientWithSecret(id, secret string) (*Client, error)
	GetRefreshToken(refreshToken string) (*RefreshToken, error)
	GetAuthorizeCode(code string) (*AuthorizeCode, error)
	GetAccessToken(accessToken string) (*AccessToken, error)
	SaveAccessToken(accessToken *AccessToken) error
	SaveRefreshToken(refreshToken *RefreshToken) error
	SaveAuthorizeCode(authCode *AuthorizeCode) error
	IsAvailableScope(scopes []string) (bool, error)
	RevokeRefreshToken(refreshToken string) error
	RevokeAccessToken(accessToken string) error
}
