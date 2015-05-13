package clover

import (
	"github.com/plimble/unik"
	"github.com/plimble/unik/uuid"
)

type codeRespType struct {
	authcodeStore    AuthCodeStore
	authcodeLifeTime int
	unik             unik.Generator
}

func NewCodeRespType(authcodeStore AuthCodeStore, authcodeLifeTime int) AuthorizeRespType {
	return &codeRespType{authcodeStore, authcodeLifeTime, uuid.NewV1()}
}

func (rt *codeRespType) Response(ad *AuthorizeData, userID string) *Response {
	ac, resp := rt.createAuthCode(ad, userID)
	if resp != nil {
		return resp
	}

	data := rt.createRespData(ac.Code, ad.State)

	return NewRespData(data).setRedirect(ad.RedirectURI, ad.respType.IsImplicit(), ad.State)
}

func (rt *codeRespType) Name() string {
	return "code"
}

func (rt *codeRespType) SupportGrant() string {
	return AUTHORIZATION_CODE
}

func (rt *codeRespType) IsImplicit() bool {
	return false
}

func (rt *codeRespType) createAuthCode(ad *AuthorizeData, userID string) (*AuthorizeCode, *Response) {
	ac := &AuthorizeCode{
		Code:        rt.generateAuthCode(),
		ClientID:    ad.Client.GetClientID(),
		UserID:      userID,
		Expires:     addSecondUnix(rt.authcodeLifeTime),
		Scope:       ad.Scope,
		RedirectURI: ad.RedirectURI,
	}

	if err := rt.authcodeStore.SetAuthorizeCode(ac); err != nil {
		return nil, errInternal(err.Error())
	}

	return ac, nil
}

func (rt *codeRespType) generateAuthCode() string {
	return rt.unik.Generate()
}

func (rt *codeRespType) createRespData(code, state string) map[string]interface{} {
	data := map[string]interface{}{
		"code": code,
	}

	if state != "" {
		data["state"] = state
	}

	return data
}
