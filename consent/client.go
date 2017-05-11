package consent

import (
	"crypto/rsa"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type Challenge struct {
	ID       string
	UserID   string
	ClientID string
	Scopes   []string
	Expired  int64
}

type RedirectRequest struct {
	Challenge string
	UserID    string
	Scopes    []string
}

type ConsentClient struct {
	privateKey *rsa.PrivateKey
}

func NewClient(privateKey *rsa.PrivateKey) *ConsentClient {
	return &ConsentClient{privateKey}
}

func (c *ConsentClient) Validate(challenge string) error {
	_, err := c.validateChallenge(challenge)

	return err
}

func (c *ConsentClient) RedirectURL(req *RedirectRequest) (string, error) {
	// challenge, err := c.validateChallenge(req.Challenge)
	// if err != nil {
	// 	return "", errors.WithStack(err)
	// }

	// now := time.Now().UTC().Truncate(time.Nanosecond)

	// token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
	// 	"aud":   ctx.Client.ID,
	// 	"scp":   strings.Join(ctx.Scopes, " "),
	// 	"exp":   now.Add(time.Minute * time.Duration(c.challengeLifeSpan)).Unix(),
	// 	"rst":   ctx.Form.Get("response_type"),
	// 	"sta":   ctx.Form.Get("state"),
	// 	"redir": ctx.Form.Get("redirect_uri"),
	// })

	// challenge, err := token.SignedString(c.privateKey)
	// if err != nil {
	// 	return "", errors.WithStack(err)
	// }

	// return fmt.Sprintf("%s&consent=%s", c.consentURL, challenge), nil
	return "", nil
}

func (c *ConsentClient) Deny(challenge string) (string, error) {
	return "", nil
}

func (c *ConsentClient) validateChallenge(challenge string) (jwt.MapClaims, error) {
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

	return claims, nil
}
