package clover

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/plimble/unik"
	"strings"
	"time"
)

type tokenJWTRespType struct {
	publicKeyStore PublicKeyStore
	base           *baseTokenRespType
	unik           unik.Generator
}

func NewJWTAccessTokenRespType(publicKeyStore PublicKeyStore, refreshTokenStore RefreshTokenStore, accessLifeTime, refreshLifeTime int) AccessTokenRespType {
	return &tokenJWTRespType{
		publicKeyStore: publicKeyStore,
		base:           newBaseTokenRespType(nil, refreshTokenStore, accessLifeTime, refreshLifeTime),
		unik:           unik.NewUUID1Base64(),
	}
}

func (rt *tokenJWTRespType) Response(td *TokenData, includeRefresh bool) *Response {
	expires := addSecondUnix(rt.base.accessLifeTime)

	key, err := rt.publicKeyStore.GetKey(td.ClientID)
	if err != nil {
		return errInternal(err.Error())
	}

	token, err := rt.encodeJWT(td.ClientID, td.UserID, td.Scope, expires, key.PrivateKey, key.Algorithm)
	if err != nil {
		return errInternal(err.Error())
	}

	rToken, resp := rt.base.createRefreshToken(td.ClientID, td.UserID, td.Scope, td.Data, rt.unik.Generate(), includeRefresh)
	if resp != nil {
		return resp
	}

	data := rt.base.createRespData(token, td.Scope, rToken, "", td.Data)

	return NewRespData(data)
}

func (rt *tokenJWTRespType) encodeJWT(clientID, userID string, scopes []string, expires int64, privateKey, algo string) (string, error) {
	token := jwt.New(getJWTAlgorithm(algo))
	token.Claims["id"] = rt.unik.Generate()
	token.Claims["iss"] = ""
	token.Claims["aud"] = clientID
	token.Claims["sub"] = userID
	token.Claims["exp"] = expires
	token.Claims["iat"] = time.Now().UTC()
	token.Claims["token_type"] = "bearer"
	token.Claims["scope"] = strings.Join(scopes, " ")

	return token.SignedString([]byte(privateKey))
}
