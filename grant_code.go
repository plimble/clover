package clover

import (
	"strings"

	"github.com/pkg/errors"
)

type AuthorizeCodeGrantType struct {
	AccessTokenLifespan  int
	RefreshTokenLifespan int
}

func (g *AuthorizeCodeGrantType) Validate(ctx *AccessTokenContext, tokenManager TokenManager) error {
	if ctx.Code == "" {
		return errors.WithStack(errCodeRequired)
	}
	if ctx.RedirectURI == "" {
		return errors.WithStack(errRedirectURIRequired)
	}

	authCode, err := tokenManager.GetAuthorizeCode(ctx.Code)
	if err != nil {
		return err
	}

	if authCode.ClientID != ctx.Client.ID {
		return errors.WithStack(errClientIDMisMatch)
	}

	if authCode.RedirectURI != ctx.RedirectURI {
		return errors.WithStack(errRedirectMisMatch)
	}

	if isExpireUnix(authCode.Expired) {
		return errors.WithStack(errCodeExpired)
	}

	ctx.Scopes = authCode.Scopes
	ctx.UserID = authCode.UserID
	ctx.AccessTokenLifespan = g.AccessTokenLifespan
	ctx.RefreshTokenLifespan = g.RefreshTokenLifespan

	return nil
}

func (g *AuthorizeCodeGrantType) Name() string {
	return "authorization_code"
}

func (g *AuthorizeCodeGrantType) CreateAccessToken(ctx *AccessTokenContext, tokenManager TokenManager) (*AccessTokenRes, error) {
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
