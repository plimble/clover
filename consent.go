package clover

import (
	"crypto/rsa"
	"fmt"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

type Challenge struct {
	ID       string
	UserID   string
	ClientID string
	Scopes   []string
	Expired  int64
}

//go:generate mockery -name Consent
type Consent interface {
	UrlWithChallenge(clientID, scope string) (*url.URL, string, error)
	ValidateChallenge(challenge string) (*Challenge, error)
}

type consent struct {
	privateKey        *rsa.PrivateKey
	challengeLifeSpan int
	consentUrl        string
}

func NewConsent(privateKey *rsa.PrivateKey, consentUrl string, challengeLifeSpan int) Consent {
	return &consent{privateKey, challengeLifeSpan, consentUrl}
}

func (c *consent) ValidateChallenge(challenge string) (*Challenge, error) {
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

	ch := &Challenge{
		ID:       claims["jti"].(string),
		UserID:   claims["sub"].(string),
		ClientID: claims["aud"].(string),
		Scopes:   strings.Fields(claims["scp"].(string)),
		Expired:  claims["exp"].(int64),
	}

	return ch, nil
}

func (c *consent) UrlWithChallenge(clientID, scope string) (*url.URL, string, error) {
	now := time.Now().UTC().Truncate(time.Nanosecond)

	id := uuid.New()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"jti": id,
		"aud": clientID,
		"scp": scope,
		"exp": now.Add(time.Minute * time.Duration(c.challengeLifeSpan)).Unix(),
	})

	challenge, err := token.SignedString(c.privateKey)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	u, err := url.Parse(c.consentUrl)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	q := u.Query()
	q.Add("consent", challenge)
	u.RawQuery = q.Encode()

	return u, id, nil
}
