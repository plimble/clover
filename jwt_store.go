package clover

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type jwtTokenStore struct {
	publicKeyStore PublicKeyStore
}

func newJWTTokenStore(publicKeyStore PublicKeyStore) *jwtTokenStore {
	return &jwtTokenStore{publicKeyStore}
}

func (s *jwtTokenStore) GetClient(id string) (Client, error) {
	return nil, nil
}

func (s *jwtTokenStore) SetAccessToken(accessToken *AccessToken) error {
	return nil
}

func (s *jwtTokenStore) GetAccessToken(accesstoken string) (*AccessToken, error) {
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
