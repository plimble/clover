package clover

import (
	"strings"

	"github.com/pkg/errors"
)

type RefreshTokenGrantType struct{}

func (g *RefreshTokenGrantType) Validate(ctx *AccessTokenContext, tokenManager TokenManager) error {
	if ctx.RefreshToken == "" {
		return errors.WithStack(errRefreshTokenRequired)
	}

	refreshToken, err := tokenManager.GetRefreshToken(ctx.RefreshToken)
	if err != nil {
		return errInvalidRefreshToken.WithCause(err)
	}

	if refreshToken.ClientID != ctx.Client.ID {
		return errors.WithStack(errClientIDMisMatch)
	}

	if isExpireUnix(refreshToken.Expired) {
		return errors.WithStack(errRefreshTokenExpired)
	}

	ctx.UserID = refreshToken.UserID
	ctx.Scopes = refreshToken.Scopes
	ctx.AccessTokenLifespan = refreshToken.AccessTokenLifespan
	ctx.RefreshTokenLifespan = refreshToken.RefreshTokenLifespan

	return nil
}

func (g *RefreshTokenGrantType) Name() string {
	return "refresh_token"
}

func (g *RefreshTokenGrantType) CreateAccessToken(ctx *AccessTokenContext, tokenManager TokenManager) (*AccessTokenRes, error) {
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

	if err = tokenManager.DeleteRefreshToken(ctx.RefreshToken); err != nil {
		return nil, err
	}

	return res, nil
}
