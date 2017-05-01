package clover

//go:generate mockery -name ClientStore -case underscore
type ClientStore interface {
}

//go:generate mockery -name AccessTokenStore -case underscore
type AccessTokenStore interface {
	// SetAccessToken(accessToken *AccessToken) error
	// GetAccessToken(at string) (*AccessToken, error)
}

//go:generate mockery -name RefreshTokenStore -case underscore
type RefreshTokenStore interface {
	// RemoveRefreshToken(rt string) error
	// SetRefreshToken(rt *RefreshToken) error
	// GetRefreshToken(rt string) (*RefreshToken, error)
}

//go:generate mockery -name AuthCodeStore -case underscore
type AuthCodeStore interface {
	// SetAuthorizeCode(ac *AuthorizeCode) error
	// GetAuthorizeCode(code string) (*AuthorizeCode, error)
}
