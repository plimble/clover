package token

import (
	"time"

	"github.com/plimble/clover/oauth2"
)

type AuthorizeCodeGrantType struct {
	AccessTokenLifespan  int
	RefreshTokenLifespan int
}

func (g *AuthorizeCodeGrantType) GrantRequest(req *TokenHandlerRequest, client *oauth2.Client, storage oauth2.Storage) (*GrantData, error) {
	code := req.Form.Get("code")
	if code == "" {
		return nil, InvalidRequest("Missing parameters: code required")
	}

	redirect := req.Form.Get("redirect_uri")
	if redirect == "" {
		return nil, InvalidRequest("Missing parameters: redirect_uri required")
	}

	authCode, err := storage.GetAuthorizeCode(code)
	if err != nil {
		if oauth2.IsNotFound(err) {
			return nil, InvalidRequest("authorize code is invalid")
		}

		return nil, err
	}

	if authCode.ClientID != req.ClientID {
		return nil, InvalidRequest("client from authorize code is mismatched with Authorization header")
	}

	if authCode.RedirectURI != redirect {
		return nil, InvalidRequest("redirect_uri from request is invalid")
	}

	if !authCode.Valid() {
		return nil, InvalidRequest("code is expired")
	}

	return &GrantData{
		Scopes:               authCode.Scopes,
		UserID:               authCode.UserID,
		AccessTokenLifespan:  g.AccessTokenLifespan,
		RefreshTokenLifespan: g.RefreshTokenLifespan,
		Extras:               authCode.Extras,
		IncludeRefreshToken:  true,
	}, nil
}

func (g *AuthorizeCodeGrantType) Name() string {
	return "authorization_code"
}

func (g *AuthorizeCodeGrantType) CreateToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, string, error) {
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

func (g *AuthorizeCodeGrantType) createAccessToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, error) {
	atoken, err := tokenGen.CreateAccessToken(&oauth2.CreateAccessTokenRequest{
		ClientID:  client.ID,
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

func (g *AuthorizeCodeGrantType) createRefreshToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, error) {
	rtoken := tokenGen.CreateRefreshToken()

	rt := &oauth2.RefreshToken{
		RefreshToken: rtoken,
		ClientID:     client.ID,
		Scopes:       grantData.Scopes,
		Expired:      time.Now().UTC().Add(time.Second * time.Duration(grantData.RefreshTokenLifespan)).Unix(),
		Extras:       grantData.Extras,
	}

	if err := storage.SaveRefreshToken(rt); err != nil {
		return "", err
	}

	return rtoken, nil
}
