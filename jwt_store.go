package clover

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type jwtAccessTokenStore struct {
	publicKeyStore PublicKeyStore
}

func newJWTAccessTokenStore(publicKeyStore PublicKeyStore) *jwtAccessTokenStore {
	return &jwtAccessTokenStore{publicKeyStore}
}

func (s *jwtAccessTokenStore) GetClient(id string) (Client, error) {
	return nil, nil
}

func (s *jwtAccessTokenStore) SetAccessToken(accessToken *AccessToken) error {
	return nil
}

func (s *jwtAccessTokenStore) GetAccessToken(accesstoken string) (*AccessToken, error) {
	token, err := jwt.Parse(accesstoken, func(token *jwt.Token) (interface{}, error) {
		key, err := s.publicKeyStore.GetKey(token.Claims["aud"].(string))
		if err != nil {
			return nil, err
		}

		return []byte(key.PublicKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	at := &AccessToken{
		AccessToken: token.Claims["id"].(string),
		ClientID:    token.Claims["aud"].(string),
		UserID:      token.Claims["sub"].(string),
		Expires:     int64(token.Claims["exp"].(float64)),
	}

	scope := token.Claims["scope"].(string)
	at.Scope = strings.Split(scope, " ")

	return at, nil
}
