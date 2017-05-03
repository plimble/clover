package clover

import (
	"strings"
)

type ClientCredentialsGrantType struct {
	tokenManager TokenManager
	config       *Config
}

func NewClientCredentialsGrantType(tokenManager TokenManager, config *Config) *ClientCredentialsGrantType {
	return &ClientCredentialsGrantType{tokenManager, config}
}

func (g *ClientCredentialsGrantType) Validate(req *AccessTokenReq) error {
	if req.Client.Public {
		return ErrPublicClient
	}

	scopes, err := DefaultGrantCheckScope(req.Scopes, req.Client.Scopes)
	if err != nil {
		return err
	}

	req.Scopes = scopes

	return nil
}

func (g *ClientCredentialsGrantType) Name() string {
	return "client_credentials"
}

func (g *ClientCredentialsGrantType) CreateAccessToken(req *AccessTokenReq) (*AccessTokenRes, error) {
	accessToken, err := g.tokenManager.GenerateAccessToken(req.Client.ID, req.UserID, g.config.AccessTokenLifespan, req.Scopes)
	if err != nil {
		return nil, err
	}

	res := &AccessTokenRes{
		AccessToken: accessToken.AccessToken,
		TokenType:   "bearer",
		ExpiresIn:   g.config.AccessTokenLifespan,
		Scope:       strings.Join(req.Scopes, " "),
		UserID:      req.UserID,
	}

	if g.config.EnableRefreshToken {
		refreshToken, err := g.tokenManager.GenerateRefreshToken(req.Client.ID, req.UserID, g.config.RefreshTokenLifespan, req.Scopes)
		if err != nil {
			return nil, err
		}

		res.RefreshToken = refreshToken.RefreshToken
	}

	return res, nil
}
