package clover

import (
	"strings"

	"github.com/pkg/errors"
)

type PassowrdGrantType struct {
	UserService          UserService
	AccessTokenLifespan  int
	RefreshTokenLifespan int
}

func (g *PassowrdGrantType) Validate(ctx *AccessTokenContext, tokenManager TokenManager) error {
	if ctx.Username == "" || ctx.Password == "" {
		return errors.WithStack(errUsernamePasswordRequired)
	}

	user, err := g.UserService.GetUser(ctx.Username, ctx.Password)
	if err != nil {
		return errInvalidUsernamePassword.WithCause(errors.WithStack(err))
	}

	ctx.Scopes = user.Scopes
	ctx.UserID = user.ID
	ctx.AccessTokenLifespan = g.AccessTokenLifespan
	ctx.RefreshTokenLifespan = g.RefreshTokenLifespan

	return nil
}

func (g *PassowrdGrantType) Name() string {
	return "password"
}

func (g *PassowrdGrantType) CreateAccessToken(ctx *AccessTokenContext, tokenManager TokenManager) (*AccessTokenRes, error) {
	accessToken, refreshToken, err := tokenManager.GenerateAccessToken(ctx, ctx.Client.IncludeRefreshToken)
	if err != nil {
		return nil, err
	}

	res := &AccessTokenRes{
		AccessToken: accessToken.AccessToken,
		TokenType:   "bearer",
		ExpiresIn:   ctx.AccessTokenLifespan,
		Scope:       strings.Join(accessToken.Scopes, " "),
		UserID:      accessToken.UserID,
	}

	if ctx.Client.IncludeRefreshToken {
		res.RefreshToken = refreshToken.RefreshToken
	}

	return res, nil
}
