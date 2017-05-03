package clover

import (
	"strings"

	"github.com/pkg/errors"
)

type RefreshGrantType struct {
	tokenManager TokenManager
	config       *Config
}

func NewRefreshGrantType(tokenManager TokenManager, config *Config) *RefreshGrantType {
	return &RefreshGrantType{tokenManager, config}
}

func (g *RefreshGrantType) Validate(req *AccessTokenReq, client *Client) error {
	token := req.Form.Get("refresh_token")
	if token == "" {
		return errors.WithStack(ErrRefreshTokenRequired)
	}

	refreshToken, err := g.tokenManager.GetRefreshToken(token)
	if err != nil {
		return err
	}

	if refreshToken.ClientID != req.Client.ID {
		return errors.WithStack(ErrClientIDMisMatch)
	}

	if isExpireUnix(refreshToken.Expired) {
		return errors.WithStack(ErrRefreshTokenExpired)
	}

	req.UserID = refreshToken.UserID
	req.Scopes = refreshToken.Scopes

	return nil
}

func (g *RefreshGrantType) Name() string {
	return "refresh_token"
}

func (g *RefreshGrantType) CreateAccessToken(req *AccessTokenReq) (*AccessTokenRes, error) {
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

	if err = g.tokenManager.DeleteRefreshToken(req.Form.Get("refresh_token")); err != nil {
		return nil, ErrInternalServer.WithCause(errors.WithStack(err))
	}

	return res, nil
}
