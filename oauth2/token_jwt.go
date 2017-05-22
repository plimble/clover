package oauth2

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
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
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (c *JWTTokenGenerator) CreateAccessToken(req *CreateAccessTokenRequest) (string, error) {
	id, err := c.genrateID()
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":         id,
		"jti":        id,
		"iss":        c.issuer,
		"aud":        req.ClientID,
		"sub":        req.UserID,
		"exp":        now.Add(time.Second * time.Duration(req.ExpiresIn)).Unix(),
		"iat":        now.Unix(),
		"token_type": "bearer",
		"scope":      strings.Join(req.Scopes, " "),
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

func (c *JWTTokenGenerator) CreateCode() string {
	return uuid.NewV4().String()
}

func (c *JWTTokenGenerator) CreateRefreshToken() string {
	return uuid.NewV4().String()
}
