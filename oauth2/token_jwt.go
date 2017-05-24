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
		"extra":      req.Extras,
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

type JWTAccessToken struct {
	Audience  string
	ExpiresAt int64
	ID        string
	IssuedAt  int64
	Issuer    string
	Subject   string
	Extra     map[string]interface{}
	Scopes    []string
}

func (a *JWTAccessToken) Valid() bool {
	return a != nil && time.Now().UTC().Unix() > a.ExpiresAt
}

func (a *JWTAccessToken) HasScope(scopes ...string) bool {
	for _, scope := range scopes {
		if ok := HierarchicScope(scope, a.Scopes); !ok {
			return false
		}
	}

	return true
}

func ClaimJWTAccessToken(publicKey *rsa.PublicKey, accesstoken string) (*JWTAccessToken, error) {
	jwttoken, err := jwt.Parse(accesstoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil || !jwttoken.Valid {
		return nil, errors.New("Invalid token")
	}

	claims, ok := jwttoken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid jwt")
	}

	at := &JWTAccessToken{
		Audience:  claims["aud"].(string),
		ExpiresAt: claims["exp"].(int64),
		ID:        claims["jti"].(string),
		IssuedAt:  claims["iat"].(int64),
		Issuer:    claims["iss"].(string),
		Subject:   claims["sub"].(string),
		Extra:     claims["extra"].(map[string]interface{}),
		Scopes:    strings.Fields(claims["scope"].(string)),
	}

	return at, nil
}
