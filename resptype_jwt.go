package clover

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/plimble/unik"
	"strings"
	"time"
)

type jwtResponseType struct {
	unik           unik.Generator
	config         *Config
	publicKeyStore PublicKeyStore
	*tokenResponseType
}

func newJWTResponseType(publicKeyStore PublicKeyStore, tokenStore AccessTokenStore, refreshStore RefreshTokenStore, config *Config) *jwtResponseType {
	return &jwtResponseType{
		unik:              unik.NewUUID1Base64(),
		config:            config,
		publicKeyStore:    publicKeyStore,
		tokenResponseType: newTokenResponseType(tokenStore, refreshStore, config),
	}
}

func (rt *jwtResponseType) createAccessToken(clientID, userID string, scopes []string) (*AccessToken, *response) {
	expires := addSecondUnix(rt.config.AccessLifeTime)

	key, err := rt.publicKeyStore.GetKey(clientID)
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

	if err := rt.tokenStore.SetAccessToken(at); err != nil {
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
