package token

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type JWTTokenGenerator struct {
	privateKey *rsa.PrivateKey
	issuer     string
}

func NewJWTTokenGenerator(privateKey *rsa.PrivateKey, issuer string) *JWTTokenGenerator {
	return &JWTTokenGenerator{privateKey, issuer}
}

func (c *JWTTokenGenerator) genrateID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", errors.WithStack(err)
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (c *JWTTokenGenerator) Generate(clientID, userID, scope string, tokenLifespan int) (string, error) {
	id, err := c.genrateID()
	if err != nil {
		return "", errors.WithStack(err)
	}

	now := time.Now().UTC().Truncate(time.Nanosecond)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":         id,
		"jti":        id,
		"iss":        c.issuer,
		"aud":        clientID,
		"sub":        userID,
		"exp":        now.Add(time.Minute * time.Duration(tokenLifespan)).Unix(),
		"iat":        now.Unix(),
		"token_type": "bearer",
		"scope":      scope,
	})

	return token.SignedString(c.privateKey)
}

func (c *JWTTokenGenerator) Validate(token string) (jwt.MapClaims, error) {
	jwttoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return &c.privateKey.PublicKey, nil
	})
	if err != nil || !jwttoken.Valid {
		return nil, errors.New("Invalid token")
	}

	claims, ok := jwttoken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid jwt")
	}

	return claims, nil
}
