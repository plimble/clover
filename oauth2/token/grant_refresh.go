package token

import (
	"time"

	"github.com/plimble/clover/oauth2"
	"go.uber.org/zap"
)

type RefreshTokenGrantType struct {
	*zap.Logger
}

func (g *RefreshTokenGrantType) GrantRequest(req *TokenHandlerRequest, client *oauth2.Client, storage oauth2.Storage) (*GrantData, error) {
	rtoken := req.Form.Get("refresh_token")
	if rtoken == "" {
		return nil, ErrRefreshTokenRequired
	}

	rt, err := storage.GetRefreshToken(rtoken)
	if err != nil {
		err = ErrUnableGetRefreshToken.WithCause(err)
		g.Error("cannot get refreshtoken",
			zap.String("refreshtoken", rtoken),
			zap.Any("error", err),
		)
		return nil, err
	}

	if rt.ClientID != req.ClientID {
		return nil, ErrClientIDMisMatched
	}

	if !rt.Valid() {
		return nil, ErrRefreshTokenExpired
	}

	return &GrantData{
		RefreshToken:         rt.RefreshToken,
		UserID:               rt.UserID,
		Scopes:               rt.Scopes,
		AccessTokenLifespan:  rt.AccessTokenLifespan,
		RefreshTokenLifespan: rt.RefreshTokenLifespan,
		Extras:               rt.Extras,
		IncludeRefreshToken:  true,
	}, nil
}

func (g *RefreshTokenGrantType) Name() string {
	return "refresh_token"
}

func (g *RefreshTokenGrantType) CreateToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, string, error) {
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

	if err = storage.RevokeRefreshToken(grantData.RefreshToken); err != nil {
		err = ErrUnableRevokeRefreshToken.WithCause(err)
		g.Error("cannot revoke refreshtoken",
			zap.String("refreshtoken", rtoken),
			zap.Any("error", err),
		)

		return "", "", err
	}

	return atoken, rtoken, err
}

func (g *RefreshTokenGrantType) createAccessToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, error) {
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

func (g *RefreshTokenGrantType) createRefreshToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, error) {
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
