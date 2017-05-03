package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ory-am/fosite"
	"github.com/pkg/errors"
)

var (
	AuthCodeEntropy = 32
	SecretLength    = 32
)

// HMAC is responsible for generating and validating challenges.
type HMAC struct {
	secret []byte
}

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

func New(secret []byte) *HMAC {
	hmac := &HMAC{
		secret: secret,
	}

	return hmac
}

// Generate generates a token and a matching signature or returns an error.
// This method implements rfc6819 Section 5.1.4.2.2: Use High Entropy for Secrets.
func (c *HMAC) Generate() (string, string, error) {
	if len(c.secret) < SecretLength/2 {
		return "", "", errors.New("Secret is not strong enough")
	}

	key, err := RandomBytes(AuthCodeEntropy)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	if len(key) < AuthCodeEntropy {
		return "", "", errors.New("Could not read enough random data for key generation")
	}

	useSecret := append([]byte{}, c.secret...)
	mac := hmac.New(sha256.New, useSecret)
	_, err = mac.Write(key)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	signature := mac.Sum([]byte{})
	encodedSignature := b64.EncodeToString(signature)
	encodedToken := fmt.Sprintf("%s.%s", b64.EncodeToString(key), encodedSignature)
	return encodedToken, encodedSignature, nil
}

// Validate validates a token and returns its signature or an error if the token is not valid.
func (c *HMAC) Validate(token string) error {
	split := strings.Split(token, ".")
	if len(split) != 2 {
		return errors.WithStack(fosite.ErrInvalidTokenFormat)
	}

	key := split[0]
	signature := split[1]
	if key == "" || signature == "" {
		return errors.WithStack(fosite.ErrInvalidTokenFormat)
	}

	decodedSignature, err := b64.DecodeString(signature)
	if err != nil {
		return errors.WithStack(err)
	}

	decodedKey, err := b64.DecodeString(key)
	if err != nil {
		return errors.WithStack(err)
	}

	useSecret := append([]byte{}, c.secret...)
	mac := hmac.New(sha256.New, useSecret)
	_, err = mac.Write(decodedKey)
	if err != nil {
		return errors.WithStack(err)
	}

	if !hmac.Equal(decodedSignature, mac.Sum([]byte{})) {
		// Hash is invalid
		return errors.WithStack(fosite.ErrTokenSignatureMismatch)
	}

	return nil
}

func (c *HMAC) Signature(token string) string {
	split := strings.Split(token, ".")

	if len(split) != 2 {
		return ""
	}

	return split[1]
}
