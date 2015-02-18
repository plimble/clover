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

func (rt *codeResponseType) GetAccessToken(td *TokenData, c *Clover, includeRefresh bool) *Response {
	return nil
}

func (rt *codeResponseType) GetAuthorizeResponse(ad *AuthorizeData, c *Clover) *Response {
	code := rt.generateAuthorizationCode()
	ac := &AuthorizeCode{
		Code:        code,
		ClientID:    ad.Client.ClientID,
		UserID:      ad.Client.UserID,
		Expires:     addSecondUnix(c.Config.AuthCodeLifetime),
		Scope:       ad.Scope,
		RedirectURI: ad.RedirectURI,
	}

	if err := c.Config.Store.SetAuthorizeCode(ac); err != nil {
		return errInternal(err.Error())
	}

	output := map[string]interface{}{
		"code": code,
	}

	if ad.State != "" {
		output["state"] = ad.State
	}

	return NewRespData(output).SetRedirect(ad)
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
