package clover

import (
	"github.com/plimble/unik"
)

type authTokenRespType struct {
	base *baseTokenRespType
	unik unik.Generator
}

func NewImplicitRespType(accessTokenstore AccessTokenStore, refreshTokenStore RefreshTokenStore, accessLifeTime, refreshLifeTime int) AuthorizeRespType {
	return &authTokenRespType{
		base: newBaseTokenRespType(accessTokenstore, refreshTokenStore, accessLifeTime, refreshLifeTime),
		unik: unik.NewUUID1Base64(),
	}
}

func (rt *authTokenRespType) Response(ad *AuthorizeData, userID string) *Response {
	at, resp := rt.base.createAccessToken(ad.Client.GetClientID(), userID, ad.Scope, ad.Client.GetData(), rt.unik.Generate())
	if resp != nil {
		return resp
	}

	data := rt.base.createRespData(at.AccessToken, at.Scope, "", ad.State, at.Data)

	return newRespData(data).setRedirect(ad.RedirectURI, ad.respType.IsImplicit(), ad.State)
}

func (rt *authTokenRespType) Name() string {
	return "token"
}

func (rt *authTokenRespType) SupportGrant() string {
	return IMPLICIT
}

func (rt *authTokenRespType) IsImplicit() bool {
	return true
}
