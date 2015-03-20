package clover

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/plimble/unik"
)

type authCodeResponseType struct {
	store  AuthServerStore
	config *Config
	unik   unik.Generator
}

func newAuthCodeResponseType(store AuthServerStore, config *Config) *authCodeResponseType {
	return &authCodeResponseType{store, config, unik.NewUUIDV1()}
}

func (rt *authCodeResponseType) GetAuthResponse(ar *authorizeRequest, client Client, scopes []string) *response {
	ac, resp := rt.createAuthCode(client, scopes, ar.redirectURI)
	if resp != nil {
		return resp
	}

	data := rt.createRespData(ac.Code, ar.state)

	return NewRespData(data).SetRedirect(ar.redirectURI, ar.responseType, ar.state)
}

func (rt *authCodeResponseType) createAuthCode(client Client, scopes []string, redirectURI string) (*AuthorizeCode, *response) {
	ac := &AuthorizeCode{
		Code:        rt.generateAuthCode(),
		ClientID:    client.GetClientID(),
		UserID:      client.GetUserID(),
		Expires:     addSecondUnix(rt.config.AuthCodeLifetime),
		Scope:       scopes,
		RedirectURI: redirectURI,
	}

	if err := rt.store.SetAuthorizeCode(ac); err != nil {
		return nil, errInternal(err.Error())
	}

	return ac, nil
}

func (rt *authCodeResponseType) generateAuthCode() string {
	code := rt.unik.Generate()
	hasher := sha512.New()
	hasher.Write([]byte(code))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (rt *authCodeResponseType) createRespData(code, state string) respData {
	data := respData{
		"code": code,
	}

	if state != "" {
		data["state"] = state
	}

	return data
}
