package clover

import (
	"strings"

	"github.com/pkg/errors"
)

type PassowrdGrantType struct {
	tokenManager TokenManager
	config       *Config
	userStore    UserStore
}

func NewPassowrdGrantType(tokenManager TokenManager, config *Config, userStore UserStore) *PassowrdGrantType {
	return &PassowrdGrantType{tokenManager, config, userStore}
}

func (g *PassowrdGrantType) Validate(req *AccessTokenReq) error {
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	if username == "" || password == "" {
		return errors.WithStack(ErrUsernamePasswordRequired)
	}

	user, err := g.userStore.GetUser(req.Form.Get("username"), req.Form.Get("password"))
	if err != nil {
		return ErrInvalidUsernamePassword.WithCause(errors.WithStack(err))
	}

	scopes, err := DefaultGrantCheckScope(req.Scopes, req.Client.Scopes)
	if err != nil {
		return err
	}

	req.UserID = user.ID
	req.Scopes = scopes

	return nil
}

func (g *PassowrdGrantType) Name() string {
	return "password"
}

func (g *PassowrdGrantType) CreateAccessToken(req *AccessTokenReq) (*AccessTokenRes, error) {
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
