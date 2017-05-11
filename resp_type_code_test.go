package clover_test

import (
	"testing"

	"github.com/plimble/clover"
	"github.com/plimble/clover/mocks"
	"github.com/stretchr/testify/require"
)

func TestCodeResponseType(t *testing.T) {
	t.Run("Name", testCodeResponseTypeName)
	t.Run("GenerateUrl", testCodeResponseTypeGenerateUrl)
}

func testCodeResponseTypeName(t *testing.T) {
	rt := &clover.CodeResponseType{
		AuthorizeCodeLifespan: 15,
	}

	require.Equal(t, "code", rt.Name())
}

func testCodeResponseTypeGenerateUrl(t *testing.T) {
	rt := &clover.CodeResponseType{
		AuthorizeCodeLifespan: 15,
	}

	authCtx := &clover.AuthorizeContext{
		RedirectURI: "http://example.com",
		State:       "state1234",
	}
	authCode := &clover.AuthorizeCode{
		Code: "code1234",
	}
	tokenManager := &mocks.TokenManager{}

	tokenManager.On("GenerateAuthorizeCode", authCtx, 15).Return(authCode, nil)

	u, err := rt.GenerateUrl(authCtx, tokenManager)
	require.NoError(t, err)
	require.Equal(t, "http://example.com?code=code1234&state=state1234", u.String())
}
