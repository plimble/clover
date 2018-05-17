package oauth2

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateJWTAccessToken(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)

	gen := NewJWTTokenGenerator(privateKey)

	token, err := gen.CreateAccessToken(&CreateAccessTokenRequest{
		ClientID:  "c1",
		UserID:    "u1",
		Scopes:    []string{"s1", "s2"},
		ExpiresIn: 3600,
		Extras: map[string]interface{}{
			"e1": "val",
			"e2": "val",
		},
	})

	require.NoError(t, err)

	jwttoken, err := ClaimJWTAccessToken(&privateKey.PublicKey, token)
	require.NoError(t, err)
	require.Equal(t, "c1", jwttoken.Audience)
	require.Equal(t, "u1", jwttoken.Subject)
	require.Equal(t, []string{"s1", "s2"}, jwttoken.Scopes)
	require.Equal(t, map[string]interface{}{
		"e1": "val",
		"e2": "val",
	}, jwttoken.Extras)
}

func TestCreateJWTAccessTokenExpired(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)

	gen := NewJWTTokenGenerator(privateKey)

	token, err := gen.CreateAccessToken(&CreateAccessTokenRequest{
		ClientID:  "c1",
		UserID:    "u1",
		Scopes:    []string{"s1", "s2"},
		ExpiresIn: -10,
		Extras: map[string]interface{}{
			"e1": "val",
			"e2": "val",
		},
	})

	require.NoError(t, err)

	jwttoken, err := ClaimJWTAccessToken(&privateKey.PublicKey, token)
	require.Error(t, err)
	require.Nil(t, jwttoken)
}

func TestJWTTokenGeneratorCreateAuthorizeCode(t *testing.T) {
	gen := NewJWTTokenGenerator(nil)

	code := gen.CreateCode()
	require.NotEmpty(t, code)
}

func TestJWTTokenGeneratorCreateRefreshToken(t *testing.T) {
	gen := NewJWTTokenGenerator(nil)

	token := gen.CreateRefreshToken()
	require.NotEmpty(t, token)
}
