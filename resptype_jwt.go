package clover

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/plimble/unik"
	"strings"
	"time"
)

type jwtResponseType struct {
	unik   unik.Generator
	config *AuthConfig
	*tokenRespType
}

func newJWTResponseType(config *AuthConfig, unik unik.Generator, tokenRespType *tokenRespType) *jwtResponseType {
	rt := &jwtResponseType{
		unik:          unik,
		config:        config,
		tokenRespType: tokenRespType,
	}

	rt.tokenRespType.createTokenFunc = rt.createToken

	return rt
}

func (rt *jwtResponseType) createToken(clientID, userID string, scopes []string, data map[string]interface{}) (*AccessToken, *Response) {
	expires := addSecondUnix(rt.config.AccessLifeTime)

	key, err := rt.config.PublicKeyStore.GetKey(clientID)
	if err != nil {
		return nil, errInternal(err.Error())
	}

	token, err := rt.encodeJWT(clientID, userID, scopes, expires, key.PrivateKey, key.Algorithm)
	if err != nil {
		return nil, errInternal(err.Error())
	}

	at := &AccessToken{
		AccessToken: token,
		ClientID:    clientID,
		UserID:      userID,
		Expires:     expires,
		Scope:       scopes,
		Data:        data,
	}

	if err := rt.config.AuthServerStore.SetAccessToken(at); err != nil {
		return nil, errInternal(err.Error())
	}

	return at, nil
}

func (rt *jwtResponseType) encodeJWT(clientID, userID string, scopes []string, expires int64, privateKey, algo string) (string, error) {
	token := jwt.New(getJWTAlgorithm(algo))
	token.Claims["id"] = rt.unik.Generate()
	token.Claims["iss"] = ""
	token.Claims["aud"] = clientID
	token.Claims["sub"] = userID
	token.Claims["exp"] = expires
	token.Claims["iat"] = time.Now()
	token.Claims["token_type"] = "bearer"
	token.Claims["scope"] = strings.Join(scopes, " ")

	return token.SignedString([]byte(privateKey))
}
