package clover

type Config struct {
	// AccessTokenLifespan sets how long an access token is going to be valid.
	AccessTokenLifespan int
	// AuthorizeCodeLifespan sets how long an authorize code is going to be valid.
	AuthorizeCodeLifespan int
	// IDTokenLifespan sets how long an id token is going to be valid.
	IDTokenLifespan int

	ScopeNotAllowedPublicClient string
}
