package clover

import (
	"github.com/plimble/unik"
	"strings"
)

type tokenResponseType struct{}

func newTokenResponseType() *tokenResponseType {
	return &tokenResponseType{}
}

func (rt *tokenResponseType) GetAuthorizeResponse(client Client, scopes []string, ar *authorizeRequest, a *AuthorizeServer) *Response {
	at, resp := rt.createAccessToken(client.GetClientID(), client.GetUserID(), scopes, a)
	if resp != nil {
		return resp
	}

	output := map[string]interface{}{
		"access_token": at.AccessToken,
		"token_type":   "bearer",
		"expires_in":   a.Config.AccessLifeTime,
		"scope":        strings.Join(at.Scope, " "),
	}

	if ar.state != "" {
		output["state"] = ar.state
	}

	return NewRespData(output).SetRedirect(ar.redirectURI, ar.responseType, ar.state)
}

func (rt *tokenResponseType) GetAccessToken(td *TokenData, a *AuthorizeServer, includeRefresh bool) *Response {
	at, resp := rt.createAccessToken(td.GrantData.ClientID, td.GrantData.UserID, td.Scope, a)
	if resp != nil {
		return resp
	}

	output := map[string]interface{}{
		"access_token": at.AccessToken,
		"token_type":   "bearer",
		"expires_in":   a.Config.AccessLifeTime,
		"scope":        strings.Join(at.Scope, " "),
	}

	if includeRefresh {
		r, resp := rt.createRefreshToken(at, a)
		if resp != nil {
			return resp
		}

		if r != nil {
			output["refresh_token"] = r.RefreshToken
		}
	}

	return NewRespData(output)
}

func (rt *tokenResponseType) createAccessToken(clientID, userID string, scope []string, a *AuthorizeServer) (*AccessToken, *Response) {
	token := unik.NewUUID1Base64().Generate()

	at := &AccessToken{
		AccessToken: token,
		ClientID:    clientID,
		UserID:      userID,
		Expires:     addSecondUnix(a.Config.AccessLifeTime),
		Scope:       scope,
	}

	if err := a.Config.Store.SetAccessToken(at); err != nil {
		return nil, errInternal(err.Error())
	}

	return at, nil
}

func (rt *tokenResponseType) createRefreshToken(at *AccessToken, a *AuthorizeServer) (*RefreshToken, *Response) {
	if a.Config.RefreshTokenLifetime < 1 {
		return nil, nil
	}

	token := unik.NewUUID1Base64().Generate()
	r := &RefreshToken{
		RefreshToken: token,
		ClientID:     at.ClientID,
		UserID:       at.UserID,
		Expires:      addSecondUnix(a.Config.AccessLifeTime),
		Scope:        at.Scope,
	}

	if err := a.Config.Store.SetRefreshToken(r); err != nil {
		return nil, errInternal(err.Error())
	}

	return r, nil
}

func (rt *tokenResponseType) GetResponseType() string {
	return "token"
}
