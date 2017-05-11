package clover

import (
	"strings"

	"github.com/pkg/errors"
)

type ClientCredentialsGrantType struct {
	AccessTokenLifespan int
}

func (g *ClientCredentialsGrantType) Validate(ctx *AccessTokenContext, tokenManager TokenManager) error {
	if ctx.Client.Public {
		return errors.WithStack(errPublicClient)
	}

	ctx.Scopes = ctx.Client.Scopes
	ctx.AccessTokenLifespan = g.AccessTokenLifespan

	return nil
}

func (g *ClientCredentialsGrantType) Name() string {
	return "client_credentials"
}

func (g *ClientCredentialsGrantType) CreateAccessToken(ctx *AccessTokenContext, tokenManager TokenManager) (*AccessTokenRes, error) {
	accessToken, _, err := tokenManager.GenerateAccessToken(ctx, false)
	if err != nil {
		return nil, err
	}

	return &AccessTokenRes{
		AccessToken: accessToken.AccessToken,
		TokenType:   "bearer",
		ExpiresIn:   ctx.AccessTokenLifespan,
		Scope:       strings.Join(accessToken.Scopes, " "),
		UserID:      accessToken.UserID,
	}, nil
}
