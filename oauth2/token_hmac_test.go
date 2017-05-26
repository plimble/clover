package oauth2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateHMACAccessToken(t *testing.T) {
	gen := NewHMACTokenGenerator([]byte("1234asdasdwfvx25hgk0dge3gfvlk"))

	token, err := gen.CreateAccessToken(nil)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	gen = NewHMACTokenGenerator([]byte("1234"))

	token, err = gen.CreateAccessToken(nil)
	require.Equal(t, "Secret is not strong enough", err.Error())
	require.Empty(t, token)
}

func TestHMACTokenGeneratorCreateAuthorizeCode(t *testing.T) {
	gen := NewHMACTokenGenerator([]byte("1234"))

	code := gen.CreateCode()
	require.NotEmpty(t, code)
}

func TestHMACTokenGeneratorCreateRefreshToken(t *testing.T) {
	gen := NewHMACTokenGenerator([]byte("1234"))

	token := gen.CreateRefreshToken()
	require.NotEmpty(t, token)
}
