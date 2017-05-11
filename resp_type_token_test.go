package clover_test

import (
	"testing"

	"github.com/plimble/clover"
	"github.com/plimble/clover/mocks"
	"github.com/stretchr/testify/require"
)

func TestTokenResponseType(t *testing.T) {
	t.Run("Name", testTokenResponseTypeName)
	t.Run("GenerateUrl", testTokenResponseTypeGenerateUrl)
}

func testTokenResponseTypeName(t *testing.T) {
	rt := &clover.TokenResponseType{
		AccessTokenLifeSpan: 15,
	}

	require.Equal(t, "token", rt.Name())
}

func testTokenResponseTypeGenerateUrl(t *testing.T) {
	rt := &clover.TokenResponseType{
		AccessTokenLifeSpan: 15,
	}

	acCtx := &clover.AuthorizeContext{
		Client:      clover.Client{ID: "c1"},
		UserID:      "u1",
		Scopes:      []string{"s1", "s2"},
		RedirectURI: "http://example.com",
		State:       "state1234",
	}
	atCtx := &clover.AccessTokenContext{
		Client:              clover.Client{ID: "c1"},
		UserID:              "u1",
		Scopes:              []string{"s1", "s2"},
		AccessTokenLifespan: 15,
	}
	at := &clover.AccessToken{
		AccessToken: "foo",
	}

	tokenManager := &mocks.TokenManager{}
	tokenManager.On("GenerateAccessToken", atCtx, false).Return(at, nil, nil)

	u, err := rt.GenerateUrl(acCtx, tokenManager)
	require.NoError(t, err)
	require.Equal(t, "http://example.com#state=state1234&token=foo", u.String())
}
