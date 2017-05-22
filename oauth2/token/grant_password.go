package token

import (
	"time"

	"github.com/plimble/clover/oauth2"
	"go.uber.org/zap"
)

type PassowrdGrantType struct {
	UserService          oauth2.UserService
	AccessTokenLifespan  int
	RefreshTokenLifespan int
	*zap.Logger
}

func (g *PassowrdGrantType) GrantRequest(req *TokenHandlerRequest, client *oauth2.Client, storage oauth2.Storage) (*GrantData, error) {
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	if username == "" || password == "" {
		return nil, ErrUsernamePasswordRequired
	}

	user, err := g.UserService.GetUser(username, password)
	if err != nil {
		err = ErrInvalidUser.WithCause(err)
		g.Error("get user",
			zap.String("username", username),
			zap.Any("error", err),
		)
		return nil, err
	}

	scopes := user.Scopes
	if scopes == nil {
		scopes = client.Scopes
	}

	return &GrantData{
		Scopes:               scopes,
		UserID:               user.ID,
		AccessTokenLifespan:  g.AccessTokenLifespan,
		RefreshTokenLifespan: g.RefreshTokenLifespan,
		IncludeRefreshToken:  true,
		Extras:               user.Extras,
	}, nil
}

func (g *PassowrdGrantType) Name() string {
	return "password"
}

func (g *PassowrdGrantType) CreateToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, string, error) {
	atoken, err := g.createAccessToken(grantData, client, storage, tokenGen)
	if err != nil {
		return "", "", err
	}

	var rtoken string
	if grantData.IncludeRefreshToken {
		if rtoken, err = g.createRefreshToken(grantData, client, storage, tokenGen); err != nil {
			return "", "", err
		}
	}

	return atoken, rtoken, nil
}

func (g *PassowrdGrantType) createAccessToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, error) {
	atoken, err := tokenGen.CreateAccessToken(&oauth2.CreateAccessTokenRequest{
		ClientID:  client.ID,
		Scopes:    grantData.Scopes,
		ExpiresIn: grantData.AccessTokenLifespan,
		Extras:    grantData.Extras,
	})
	if err != nil {
		err = ErrUnableCreateAccessToken.WithCause(err)
		g.Error("cannot create accesstoken",
			zap.NamedError("cause", err),
			zap.Any("error", err),
		)
		return "", err
	}

	at := &oauth2.AccessToken{
		AccessToken: atoken,
		ClientID:    client.ID,
		Scopes:      grantData.Scopes,
		Expired:     time.Now().UTC().Add(time.Second * time.Duration(grantData.AccessTokenLifespan)).Unix(),
		ExpiresIn:   grantData.AccessTokenLifespan,
		Extras:      grantData.Extras,
	}

	if err = storage.SaveAccessToken(at); err != nil {
		err = ErrUnableCreateAccessToken.WithCause(err)
		g.Error("cannot save accesstoken",
			zap.Any("AccessToken", at),
			zap.Any("error", err),
		)
		return "", err
	}

	return atoken, nil
}

func (g *PassowrdGrantType) createRefreshToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, error) {
	rtoken := tokenGen.CreateRefreshToken()

	rt := &oauth2.RefreshToken{
		RefreshToken: rtoken,
		ClientID:     client.ID,
		Scopes:       grantData.Scopes,
		Expired:      time.Now().UTC().Add(time.Second * time.Duration(grantData.RefreshTokenLifespan)).Unix(),
		Extras:       grantData.Extras,
	}

	if err := storage.SaveRefreshToken(rt); err != nil {
		err = ErrUnableCreateRefreshToken.WithCause(err)
		g.Error("cannot save accesstoken",
			zap.Any("RefreshToken", rt),
			zap.Any("error", err),
		)
		return "", err
	}

	return rtoken, nil
}
