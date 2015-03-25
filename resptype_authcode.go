package clover

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/plimble/unik"
)

type codeRespType struct {
	config *AuthConfig
	unik   unik.Generator
}

func newCodeRespType(config *AuthConfig, unik unik.Generator) *codeRespType {
	return &codeRespType{config, unik}
}

func (rt *codeRespType) GetAuthResponse(ar *authorizeRequest, client Client, scopes []string) *Response {
	ac, resp := rt.createAuthCode(client, scopes, ar.redirectURI)
	if resp != nil {
		return resp
	}

	data := rt.createRespData(ac.Code, ar.state)

	return newRespData(data).setRedirect(ar.redirectURI, ar.responseType, ar.state)
}

func (rt *codeRespType) createAuthCode(client Client, scopes []string, redirectURI string) (*AuthorizeCode, *Response) {
	ac := &AuthorizeCode{
		Code:        rt.generateAuthCode(),
		ClientID:    client.GetClientID(),
		UserID:      client.GetUserID(),
		Expires:     addSecondUnix(rt.config.AuthCodeLifetime),
		Scope:       scopes,
		RedirectURI: redirectURI,
	}

	if err := rt.config.AuthCodeStore.SetAuthorizeCode(ac); err != nil {
		return nil, errInternal(err.Error())
	}

	return ac, nil
}

func (rt *codeRespType) generateAuthCode() string {
	code := rt.unik.Generate()
	hasher := sha512.New()
	hasher.Write([]byte(code))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (rt *codeRespType) createRespData(code, state string) respData {
	data := respData{
		"code": code,
	}

	if state != "" {
		data["state"] = state
	}

	return data
}
