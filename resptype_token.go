package clover

import (
	"github.com/plimble/unik"
	"strings"
)

type createTokenFunc func(clientID, userID string, scopes []string) (*AccessToken, *Response)

type tokenRespType struct {
	config          *AuthConfig
	unik            unik.Generator
	createTokenFunc createTokenFunc
}

func newTokenRespType(config *AuthConfig, unik unik.Generator) *tokenRespType {
	rt := &tokenRespType{
		config: config,
		unik:   unik,
	}

	rt.createTokenFunc = rt.createToken

	return rt
}

func (rt *tokenRespType) GetAuthResponse(ar *authorizeRequest, client Client, scopes []string) *Response {
	at, resp := rt.createTokenFunc(client.GetClientID(), client.GetUserID(), scopes)
	if resp != nil {
		return resp
	}

	data := rt.createRespData(at.AccessToken, rt.config.AuthCodeLifetime, scopes, "", ar.state)

	return newRespData(data).setRedirect(ar.redirectURI, ar.responseType, ar.state)
}

func (rt *tokenRespType) GetAccessToken(td *TokenData, includeRefresh bool) *Response {
	at, resp := rt.createTokenFunc(td.GrantData.ClientID, td.GrantData.UserID, td.Scope)
	if resp != nil {
		return resp
	}

	refresh, resp := rt.createRefreshToken(at, includeRefresh)
	if resp != nil {
		return resp
	}

	data := rt.createRespData(at.AccessToken, rt.config.AccessLifeTime, at.Scope, refresh, "")

	return newRespData(data)
}

func (rt *tokenRespType) createToken(clientID, userID string, scopes []string) (*AccessToken, *Response) {
	at := &AccessToken{
		AccessToken: rt.generateToken(),
		ClientID:    clientID,
		UserID:      userID,
		Expires:     addSecondUnix(rt.config.AccessLifeTime),
		Scope:       scopes,
	}

	if err := rt.config.AuthServerStore.SetAccessToken(at); err != nil {
		return nil, errInternal(err.Error())
	}

	return at, nil
}

func (rt *tokenRespType) createRefreshToken(at *AccessToken, includeRefresh bool) (string, *Response) {
	if !includeRefresh || rt.config.RefreshTokenStore == nil || rt.config.RefreshTokenLifetime < 1 {
		return "", nil
	}

	r := &RefreshToken{
		RefreshToken: rt.generateToken(),
		ClientID:     at.ClientID,
		UserID:       at.UserID,
		Expires:      addSecondUnix(rt.config.RefreshTokenLifetime),
		Scope:        at.Scope,
	}

	if err := rt.config.RefreshTokenStore.SetRefreshToken(r); err != nil {
		return "", errInternal(err.Error())
	}

	return r.RefreshToken, nil
}

func (rt *tokenRespType) generateToken() string {
	return rt.unik.Generate()
}

func (rt *tokenRespType) createRespData(token string, expiresIn int, scopes []string, refresh, state string) respData {
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
