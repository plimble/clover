package clover

import (
	"github.com/plimble/unik"
	"strings"
)

type tokenResponseType struct {
	tokenStore   AccessTokenStore
	refreshStore RefreshTokenStore
	config       *Config
	unik         unik.Generator
}

func newTokenResponseType(tokenStore AccessTokenStore, refreshStore RefreshTokenStore, config *Config) *tokenResponseType {
	return &tokenResponseType{tokenStore, refreshStore, config, unik.NewUUID1Base64()}
}

func (rt *tokenResponseType) GetAuthResponse(ar *authorizeRequest, client Client, scopes []string) *response {
	at, resp := rt.createAccessToken(client.GetClientID(), client.GetUserID(), scopes)
	if resp != nil {
		return resp
	}

	data := rt.createRespData(at.AccessToken, rt.config.AuthCodeLifetime, scopes, "", ar.state)

	return NewRespData(data).SetRedirect(ar.redirectURI, ar.responseType, ar.state)
}

func (rt *tokenResponseType) GetAccessToken(td *TokenData, includeRefresh bool) *response {
	at, resp := rt.createAccessToken(td.GrantData.ClientID, td.GrantData.UserID, td.Scope)
	if resp != nil {
		return resp
	}

	refresh, resp := rt.createRefreshToken(at, includeRefresh)
	if resp != nil {
		return resp
	}

	data := rt.createRespData(at.AccessToken, rt.config.AccessLifeTime, at.Scope, refresh, "")

	return NewRespData(data)
}

func (rt *tokenResponseType) createAccessToken(clientID, userID string, scopes []string) (*AccessToken, *response) {
	at := &AccessToken{
		AccessToken: rt.generateToken(),
		ClientID:    clientID,
		UserID:      userID,
		Expires:     addSecondUnix(rt.config.AccessLifeTime),
		Scope:       scopes,
	}

	if err := rt.tokenStore.SetAccessToken(at); err != nil {
		return nil, errInternal(err.Error())
	}

	return at, nil
}

func (rt *tokenResponseType) createRefreshToken(at *AccessToken, includeRefresh bool) (string, *response) {
	if !includeRefresh || rt.refreshStore == nil || rt.config.RefreshTokenLifetime < 1 {
		return "", nil
	}

	r := &RefreshToken{
		RefreshToken: rt.generateToken(),
		ClientID:     at.ClientID,
		UserID:       at.UserID,
		Expires:      addSecondUnix(rt.config.RefreshTokenLifetime),
		Scope:        at.Scope,
	}

	if err := rt.refreshStore.SetRefreshToken(r); err != nil {
		return "", errInternal(err.Error())
	}

	return r.RefreshToken, nil
}

func (rt *tokenResponseType) generateToken() string {
	return rt.unik.Generate()
}

func (rt *tokenResponseType) SetRefreshStore(store RefreshTokenStore) {
	rt.refreshStore = store
}

func (rt *tokenResponseType) createRespData(token string, expiresIn int, scopes []string, refresh, state string) respData {
	data := respData{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   rt.config.AccessLifeTime,
		"scope":        strings.Join(scopes, " "),
	}

	if refresh != "" {
		data["refresh_token"] = refresh
	}

	if state != "" {
		data["state"] = state
	}

	return data
}
