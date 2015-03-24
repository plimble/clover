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

func newJWTResponseType(config *AuthConfig) *jwtResponseType {
	return &jwtResponseType{
		unik:          unik.NewUUID1Base64(),
		config:        config,
		tokenRespType: newTokenRespType(config),
	}
}

func (rt *jwtResponseType) GetAuthResponse(ar *authorizeRequest, client Client, scopes []string) *Response {
	at, resp := rt.createToken(client.GetClientID(), client.GetUserID(), scopes)
	if resp != nil {
		return resp
	}

	data := rt.createRespData(at.AccessToken, rt.config.AuthCodeLifetime, scopes, "", ar.state)

	return newRespData(data).setRedirect(ar.redirectURI, ar.responseType, ar.state)
}

func (rt *jwtResponseType) GetAccessToken(td *TokenData, includeRefresh bool) *Response {
	at, resp := rt.createToken(td.GrantData.ClientID, td.GrantData.UserID, td.Scope)
	if resp != nil {
		return resp
	}

	refresh, resp := rt.createRefreshToken(at, includeRefresh)
	if resp != nil {
		return resp
	}

	data := rt.createRespData(at.AccessToken, rt.config.AccessLifeTime, at.Scope, refresh, "")

	return newRespData(data)
}

func (rt *jwtResponseType) createToken(clientID, userID string, scopes []string) (*AccessToken, *Response) {
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
