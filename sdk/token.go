package sdk

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/plimble/clover/scope"
	"golang.org/x/net/context"
)

type accessTokenKey struct{}

type AccessToken struct {
	AccessToken string
	ClientID    string
	UserID      string
	Scopes      []string
	Expired     int64
}

func (a *AccessToken) CheckScope(requestScope string) bool {
	return scope.HierarchicScope(requestScope, a.Scopes)
}

type JWTTokenValidator struct {
	publicKey *rsa.PublicKey
}

func NewJWTTokenValidator(publicKey *rsa.PublicKey) *JWTTokenValidator {
	return &JWTTokenValidator{publicKey}
}

func NewAccessTokenContext(ctx context.Context, at *AccessToken) context.Context {
	return context.WithValue(ctx, accessTokenKey{}, at)
}

func GetAccessTokenFromContext(ctx context.Context) (*AccessToken, bool) {
	at, ok := ctx.Value(accessTokenKey{}).(*AccessToken)
	return at, ok
}

func (v *JWTTokenValidator) Validate(token string, requestScope string) (*AccessToken, error) {
	jwttoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return &v.publicKey, nil
	})
	if err != nil || !jwttoken.Valid {
		return nil, errors.New("Invalid token")
	}

	claims, ok := jwttoken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid jwt")
	}

	at := &AccessToken{
		AccessToken: claims["jti"].(string),
		Scopes:      strings.Fields(claims["scope"].(string)),
		Expired:     claims["exp"].(int64),
		ClientID:    claims["aud"].(string),
		UserID:      claims["sub"].(string),
	}

	if time.Now().UTC().Truncate(time.Nanosecond).Unix() > at.Expired {
		return nil, errors.New("token expired")
	}

	return at, nil
}
