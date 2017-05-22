package token

import (
	"time"

	"github.com/plimble/clover/oauth2"
	"go.uber.org/zap"
)

type AuthorizeCodeGrantType struct {
	AccessTokenLifespan  int
	RefreshTokenLifespan int
	*zap.Logger
}

func (g *AuthorizeCodeGrantType) GrantRequest(req *TokenHandlerRequest, client *oauth2.Client, storage oauth2.Storage) (*GrantData, error) {
	code := req.Form.Get("code")
	if code == "" {
		return nil, ErrCodeRequired
	}

	redirect := req.Form.Get("redirect_uri")
	if redirect == "" {
		return nil, ErrRedirectRequired
	}

	authCode, err := storage.GetAuthorizeCode(code)
	if err != nil {
		err = ErrUnableGetAuthorizeCode.WithCause(err)
		g.Error("GetAuthorizeCode",
			zap.String("code", code),
			zap.Any("error", err),
		)

		return nil, err
	}

	if authCode.ClientID != req.ClientID {
		return nil, ErrClientIDMisMatched
	}

	if authCode.RedirectURI != redirect {
		return nil, ErrRedirectMisMatched
	}

	if !authCode.Valid() {
		return nil, ErrCodeExpired
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
		err = ErrUnableCreateRefreshToken.WithCause(err)
		g.Error("cannot save accesstoken",
			zap.Any("RefreshToken", rt),
			zap.Any("error", err),
		)
		return "", err
	}

	return rtoken, nil
}
