package clover

import (
	"strings"

	"github.com/pkg/errors"
)

type AuthorizeCodeGrantType struct {
	tokenManager TokenManager
	config       *Config
}

func NewAuthorizeCodeGrantType(tokenManager TokenManager, config *Config) *AuthorizeCodeGrantType {
	return &AuthorizeCodeGrantType{tokenManager, config}
}

func (g *AuthorizeCodeGrantType) Validate(req *AccessTokenReq) error {
	code := req.Form.Get("code")
	if code == "" {
		return errors.WithStack(ErrCodeRequired)
	}

	authCode, err := g.tokenManager.GetAuthorizeCode(code)
	if err != nil {
		return err
	}

	if authCode.ClientID != req.Client.ID {
		return errors.WithStack(ErrClientIDMisMatch)
	}

	if authCode.RedirectURI != req.Form.Get("redirect_uri") {
		return errors.WithStack(ErrRedirectMisMatch)
	}

	if isExpireUnix(authCode.Expired) {
		return errors.WithStack(ErrCodeExpired)
	}

	req.Scopes = authCode.Scopes
	req.UserID = authCode.UserID

	return nil
}

func (g *AuthorizeCodeGrantType) Name() string {
	return "client_credentials"
}

func (g *AuthorizeCodeGrantType) CreateAccessToken(req *AccessTokenReq) (*AccessTokenRes, error) {
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
