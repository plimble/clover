package oauth2

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateHMACAccessToken(t *testing.T) {
	gen := NewHMACTokenGenerator([]byte("1234asdasdwfvx25hgk0dge3gfvlk"), zap.L())

	token, err := gen.CreateAccessToken(nil)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	gen = NewHMACTokenGenerator([]byte("1234"), zap.L())

	token, err = gen.CreateAccessToken(nil)
	require.Equal(t, "Secret is not strong enough", err.Error())
	require.Empty(t, token)
}

func TestHMACTokenGeneratorCreateAuthorizeCode(t *testing.T) {
	gen := NewHMACTokenGenerator([]byte("1234"), zap.L())

	code := gen.CreateCode()
	require.NotEmpty(t, code)
}

func TestHMACTokenGeneratorCreateRefreshToken(t *testing.T) {
	gen := NewHMACTokenGenerator([]byte("1234"), zap.L())

	token := gen.CreateRefreshToken()
	require.NotEmpty(t, token)
}
