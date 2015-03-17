package clover

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/plimble/unik"
)

type codeResponseType struct{}

func newCodeResponseType() *codeResponseType {
	return &codeResponseType{}
}

func (rt *codeResponseType) GetAccessToken(td *TokenData, a *AuthorizeServer, includeRefresh bool) *Response {
	return nil
}

func (rt *codeResponseType) GetAuthorizeResponse(client Client, scopes []string, ar *authorizeRequest, a *AuthorizeServer) *Response {
	code := rt.generateAuthorizationCode()
	ac := &AuthorizeCode{
		Code:        code,
		ClientID:    client.GetClientID(),
		UserID:      client.GetUserID(),
		Expires:     addSecondUnix(a.Config.AuthCodeLifetime),
		Scope:       scopes,
		RedirectURI: ar.redirectURI,
	}

	if err := a.Config.Store.SetAuthorizeCode(ac); err != nil {
		return errInternal(err.Error())
	}

	output := map[string]interface{}{
		"code": code,
	}

	if ar.state != "" {
		output["state"] = ar.state
	}

	return NewRespData(output).SetRedirect(ar.redirectURI, ar.responseType, ar.state)
}

func (rt *codeResponseType) generateAuthorizationCode() string {
	code := unik.NewUUIDV1().Generate()
	hasher := sha512.New()
	hasher.Write([]byte(code))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (rt *codeResponseType) GetResponseType() string {
	return "code"
}
