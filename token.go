package clover

import (
	"github.com/plimble/unik"
	"strings"
)

type tokenResponseType struct{}

func newTokenResponseType() *tokenResponseType {
	return &tokenResponseType{}
}

func (rt *tokenResponseType) GetAuthorizeResponse(ad *AuthorizeData, c *Clover) *Response {
	at, resp := rt.createAccessToken(ad.Client.ClientID, ad.Client.UserID, ad.Scope, c)
	if resp != nil {
		return resp
	}

	output := map[string]interface{}{
		"access_token": at.AccessToken,
		"token_type":   "bearer",
		"expires_in":   c.Config.AccessLifeTime,
		"scope":        strings.Join(at.Scope, " "),
	}

	if ad.State != "" {
		output["state"] = ad.State
	}

	return NewRespData(output).SetRedirect(ad)
}

func (rt *tokenResponseType) GetAccessToken(td *TokenData, c *Clover, includeRefresh bool) *Response {
	at, resp := rt.createAccessToken(td.Client.ClientID, td.Client.UserID, td.Scope, c)
	if resp != nil {
		return resp
	}

	output := map[string]interface{}{
		"access_token": at.AccessToken,
		"token_type":   "bearer",
		"expires_in":   c.Config.AccessLifeTime,
		"scope":        strings.Join(at.Scope, " "),
	}

	if includeRefresh {
		r, resp := rt.createRefreshToken(at, c)
		if resp != nil {
			return resp
		}

		if r != nil {
			output["refresh_token"] = r.RefreshToken
		}
	}

	return NewRespData(output)
}

func (rt *tokenResponseType) createAccessToken(clientID, userID string, scope []string, c *Clover) (*AccessToken, *Response) {
	token := unik.NewUUID1Base64().Generate()

	at := &AccessToken{
		AccessToken: token,
		ClientID:    clientID,
		UserID:      userID,
		Expires:     addSecondUnix(c.Config.AccessLifeTime),
		Scope:       scope,
	}

	if err := c.Config.Store.SetAccessToken(at); err != nil {
		return nil, errInternal(err.Error())
	}

	return at, nil
}

func (rt *tokenResponseType) createRefreshToken(at *AccessToken, c *Clover) (*RefreshToken, *Response) {
	if c.Config.RefreshTokenLifetime < 1 {
		return nil, nil
	}

	token := unik.NewUUID1Base64().Generate()
	r := &RefreshToken{
		RefreshToken: token,
		ClientID:     at.ClientID,
		UserID:       at.UserID,
		Expires:      addSecondUnix(c.Config.AccessLifeTime),
		Scope:        at.Scope,
	}

	if err := c.Config.Store.SetRefreshToken(r); err != nil {
		return nil, errInternal(err.Error())
	}

	return r, nil
}

func (rt *tokenResponseType) GetResponseType() string {
	return "token"
}
