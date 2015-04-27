package clover

import (
	"github.com/plimble/unik"
)

type tokenAccessTokenRespType struct {
	base *baseTokenRespType
	unik unik.Generator
}

func NewAccessTokenRespType(accessTokenstore AccessTokenStore, refreshTokenStore RefreshTokenStore, accessLifeTime, refreshLifeTime int) AccessTokenRespType {
	return &tokenAccessTokenRespType{
		base: newBaseTokenRespType(accessTokenstore, refreshTokenStore, accessLifeTime, refreshLifeTime),
		unik: unik.NewUUID1Base64(),
	}
}

func (rt *tokenAccessTokenRespType) Response(td *TokenData, includeRefresh bool) *Response {
	aToken, resp := rt.base.createAccessToken(td.ClientID, td.UserID, td.Scope, td.Data, rt.unik.Generate())
	if resp != nil {
		return resp
	}

	rToken, resp := rt.base.createRefreshToken(td.ClientID, td.UserID, td.Scope, td.Data, rt.unik.Generate(), includeRefresh)
	if resp != nil {
		return resp
	}

	data := rt.base.createRespData(aToken.AccessToken, aToken.Scope, rToken, "", aToken.Data)

	return newRespData(data)
}
