package oauth2

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var (
	AuthCodeEntropy = 32
	SecretLength    = 32
)

// HMACTokenGenerator is responsible for generating and validating challenges.
type HMACTokenGenerator struct {
	secret []byte
}

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

func NewHMACTokenGenerator(secret []byte) *HMACTokenGenerator {
	hmac := &HMACTokenGenerator{
		secret: secret,
	}

	return hmac
}

// Generate generates a token and a matching signature or returns an error.
// This method implements rfc6819 Section 5.1.4.2.2: Use High Entropy for Secrets.
func (c *HMACTokenGenerator) CreateAccessToken(req *CreateAccessTokenRequest) (string, error) {
	if len(c.secret) < SecretLength/2 {
		return "", errors.New("Secret is not strong enough")
	}

	key, err := RandomBytes(AuthCodeEntropy)
	if err != nil {
		return "", err
	}

	if len(key) < AuthCodeEntropy {
		return "", errors.New("Could not read enough random data for key generation")
	}

	useSecret := append([]byte{}, c.secret...)
	mac := hmac.New(sha256.New, useSecret)
	_, err = mac.Write(key)
	if err != nil {
		return "", err
	}

	signature := mac.Sum([]byte{})
	encodedSignature := b64.EncodeToString(signature)
	encodedToken := fmt.Sprintf("%s.%s", b64.EncodeToString(key), encodedSignature)
	return encodedToken, nil
}

// Validate validates a token and returns its signature or an error if the token is not valid.
func (c *HMACTokenGenerator) Validate(token string) error {
	split := strings.Split(token, ".")
	if len(split) != 2 {
		return errors.New("invalid token format")
	}

	key := split[0]
	signature := split[1]
	if key == "" || signature == "" {
		return errors.New("invalid token format")
	}

	decodedSignature, err := b64.DecodeString(signature)
	if err != nil {
		return err
	}

	decodedKey, err := b64.DecodeString(key)
	if err != nil {
		return err
	}

	useSecret := append([]byte{}, c.secret...)
	mac := hmac.New(sha256.New, useSecret)
	_, err = mac.Write(decodedKey)
	if err != nil {
		return err
	}

	if !hmac.Equal(decodedSignature, mac.Sum([]byte{})) {
		// Hash is invalid
		return errors.New("token signature mismatch")
	}

	return nil
}

func (c *HMACTokenGenerator) Signature(token string) string {
	split := strings.Split(token, ".")

	if len(split) != 2 {
		return ""
	}

	return split[1]
}

func (c *HMACTokenGenerator) CreateCode() string {
	return uuid.NewV4().String()
}

func (c *HMACTokenGenerator) CreateRefreshToken() string {
	return uuid.NewV4().String()
}

// RandomBytes returns n random bytes by reading from crypto/rand.Reader
func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
