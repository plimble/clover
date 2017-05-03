package clover

type Config struct {
	// AccessTokenLifespan sets how long an access token is going to be valid.
	AccessTokenLifespan int
	// AuthorizeCodeLifespan sets how long an authorize code is going to be valid.
	AuthorizeCodeLifespan int
	// RefreshTokenLifespan sets how long an refresh token is going to be valid.
	RefreshTokenLifespan int

	EnableRefreshToken bool

	// hmac or jwt
	Token         string
	HMACGlobalKey string
}
