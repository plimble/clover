package sdk

import (
	"crypto/rsa"
	"fmt"
	"time"

	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type ChallengeClaim struct {
	ID          string
	ClientID    string
	Scopes      []string
	Expired     int64
	RedirectURI string
}

type RedirectRequest struct {
	Challenge string
	UserID    string
}

type Consent struct {
	privateKey *rsa.PrivateKey
}

func NewConsent(privateKey *rsa.PrivateKey) *Consent {
	return &Consent{privateKey}
}

func (c *Consent) RedirectURL(req *RedirectRequest) (string, error) {
	challenge, err := c.ValidateChallenge(req.Challenge)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"jti": challenge.ID,
		"sub": req.UserID,
		"aud": challenge.ClientID,
		"scp": strings.Join(challenge.Scopes, " "),
		"exp": challenge.Expired,
	})

	newChallenge, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return fmt.Sprintf("%s?consent=%s", challenge.RedirectURI, newChallenge), nil
}

func (c *Consent) Deny(challenge string) (string, error) {
	cc, err := c.ValidateChallenge(challenge)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s?consent=denied", cc.RedirectURI), nil
}

func (c *Consent) ValidateChallenge(challenge string) (*ChallengeClaim, error) {
	jwttoken, err := jwt.Parse(challenge, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return &c.privateKey.PublicKey, nil
	})
	if err != nil || !jwttoken.Valid {
		return nil, errors.New("Invalid challenge")
	}

	claims, ok := jwttoken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid challenge")
	}

	cc := &ChallengeClaim{
		ID:          claims["jti"].(string),
		ClientID:    claims["aud"].(string),
		Scopes:      strings.Fields(claims["scp"].(string)),
		Expired:     claims["exp"].(int64),
		RedirectURI: claims["redir"].(string),
	}

	if time.Now().UTC().Truncate(time.Nanosecond).Unix() > cc.Expired {
		return nil, errors.New("challenge expired")
	}

	return cc, nil
}
