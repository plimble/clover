package token

import (
	"time"

	"github.com/plimble/clover/oauth2"
)

type PassowrdGrantType struct {
	UserService          oauth2.UserService
	AccessTokenLifespan  int
	RefreshTokenLifespan int
}

func (g *PassowrdGrantType) GrantRequest(req *TokenHandlerRequest, client *oauth2.Client, storage oauth2.Storage) (*GrantData, error) {
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	if username == "" || password == "" {
		return nil, InvalidRequest("username and password is required")
	}

	user, err := g.UserService.GetUser(username, password)
	if err != nil {
		if oauth2.IsNotFound(err) {
			return nil, InvalidRequest("username or password is invalid")
		}

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
		UserID:    grantData.UserID,
		Scopes:    grantData.Scopes,
		ExpiresIn: grantData.AccessTokenLifespan,
		Extras:    grantData.Extras,
	})
	if err != nil {
		return "", err
	}

	at := &oauth2.AccessToken{
		AccessToken: atoken,
		ClientID:    client.ID,
		UserID:      grantData.UserID,
		Scopes:      grantData.Scopes,
		Expired:     time.Now().UTC().Add(time.Second * time.Duration(grantData.AccessTokenLifespan)).Unix(),
		ExpiresIn:   grantData.AccessTokenLifespan,
		Extras:      grantData.Extras,
	}

	if err = storage.SaveAccessToken(at); err != nil {
		return "", err
	}

	return atoken, nil
}

func (g *PassowrdGrantType) createRefreshToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, error) {
	rtoken := tokenGen.CreateRefreshToken()

	rt := &oauth2.RefreshToken{
		RefreshToken:         rtoken,
		ClientID:             client.ID,
		UserID:               grantData.UserID,
		Scopes:               grantData.Scopes,
		Expired:              time.Now().UTC().Add(time.Second * time.Duration(grantData.RefreshTokenLifespan)).Unix(),
		Extras:               grantData.Extras,
		AccessTokenLifespan:  grantData.AccessTokenLifespan,
		RefreshTokenLifespan: grantData.RefreshTokenLifespan,
	}

	if err := storage.SaveRefreshToken(rt); err != nil {
		return "", err
	}

	return rtoken, nil
}
