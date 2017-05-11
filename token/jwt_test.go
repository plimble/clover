package token

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJWTToken(t *testing.T) {
	t.Run("GenerateJWTAndValidate", testGenerateJWTAndValidate)
	t.Run("ValidateJWTFailed", testValidateJWTFailed)
}

func testGenerateJWTAndValidate(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	require.NoError(t, err)
	tokenGenerator := NewJWTTokenGenerator(privateKey, "tester")

	token, err := tokenGenerator.Generate("c1", "u1", "scope", 3600)
	require.Nil(t, err, "%s", err)
	require.NotNil(t, token)

	claims, err := tokenGenerator.Validate(token)
	require.Nil(t, err, "%s", err)
	require.NotNil(t, claims)
	require.Equal(t, "c1", claims["aud"])
	require.Equal(t, "u1", claims["sub"])
}

func testValidateJWTFailed(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	require.NoError(t, err)
	tokenGenerator := NewJWTTokenGenerator(privateKey, "tester")

	token, err := tokenGenerator.Generate("c1", "u1", "scope", 3600)
	require.Nil(t, err, "%s", err)
	require.NotNil(t, token)

	claims, err := tokenGenerator.Validate(token + "1")
	require.Equal(t, "Invalid token", err.Error())
	require.Nil(t, claims)
}
