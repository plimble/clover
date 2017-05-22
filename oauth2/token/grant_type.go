package token

import "github.com/plimble/clover/oauth2"

type GrantData struct {
	RefreshToken         string
	UserID               string
	Scopes               []string
	AccessTokenLifespan  int
	RefreshTokenLifespan int
	Extras               map[string]interface{}
	IncludeRefreshToken  bool
}

// GrantType use this interface for implement extension grant type
type GrantType interface {
	GrantRequest(req *TokenHandlerRequest, client *oauth2.Client, storage oauth2.Storage) (*GrantData, error)
	Name() string
	// CreateAccessToken return acesstoken, refreshtoken and error
	CreateToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, string, error)
}
