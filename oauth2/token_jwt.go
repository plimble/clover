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
	"go.uber.org/zap"
)

type JWTTokenGenerator struct {
	privateKey *rsa.PrivateKey
}

func NewJWTTokenGenerator(privateKey *rsa.PrivateKey) *JWTTokenGenerator {
	return &JWTTokenGenerator{privateKey}
}

func (c *JWTTokenGenerator) genrateID() (string, error) {
	bytes := make([]byte, 10)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		Logger.Error("unable to generate jwt id", zap.Error(err))

		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (c *JWTTokenGenerator) CreateAccessToken(req *CreateAccessTokenRequest) (string, error) {
	id, err := c.genrateID()
	if err != nil {
		return "", ServerError("unable to generate jwt id", err)
	}

	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"jti":        id,
		"aud":        req.ClientID,
		"sub":        req.UserID,
		"exp":        now.Add(time.Second * time.Duration(req.ExpiresIn)).Unix(),
		"iat":        now.Unix(),
		"token_type": "bearer",
		"scope":      strings.Join(req.Scopes, " "),
		"extra":      req.Extras,
	})

	atoken, err := token.SignedString(c.privateKey)
	if err != nil {
		Logger.Error("could not signed jwt string", zap.Error(err))
		return "", ServerError("could not signed jwt string", err)
	}
	return atoken, nil
}

func (c *JWTTokenGenerator) CreateCode() string {
	return uuid.NewV4().String()
}

func (c *JWTTokenGenerator) CreateRefreshToken() string {
	return uuid.NewV4().String()
}

func ClaimJWTAccessToken(publicKey *rsa.PublicKey, accesstoken string) (*JWTAccessToken, error) {
	jwttoken, err := jwt.Parse(accesstoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})

	if err != nil || !jwttoken.Valid {
		return nil, errors.New("Invalid token: " + err.Error())
	}

	claims := jwttoken.Claims.(jwt.MapClaims)

	at := &JWTAccessToken{
		Audience:  claims["aud"].(string),
		ExpiresAt: int64(claims["exp"].(float64)),
		ID:        claims["jti"].(string),
		IssuedAt:  int64(claims["iat"].(float64)),
		Subject:   claims["sub"].(string),
		Scopes:    strings.Fields(claims["scope"].(string)),
	}

	if extras, ok := claims["extra"].(map[string]interface{}); ok {
		at.Extras = extras
	}

	return at, nil
}
