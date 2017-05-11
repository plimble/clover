package clover_test

import (
	"errors"
	"testing"

	"github.com/plimble/clover"
	"github.com/plimble/clover/mocks"
	"github.com/stretchr/testify/require"
)

func TestClientGrantType(t *testing.T) {
	t.Run("Name", testClientGrantTypeName)
	t.Run("Validate", testClientGrantTypeValidate)
	t.Run("CreateAccessToken", testClientGrantTypeCreateAccessToken)
}

func testClientGrantTypeName(t *testing.T) {
	g := clover.ClientCredentialsGrantType{}
	require.Equal(t, "client_credentials", g.Name())
}

func testClientGrantTypeValidate(t *testing.T) {
	g := clover.ClientCredentialsGrantType{
		AccessTokenLifespan: 15,
	}

	atCtx := &clover.AccessTokenContext{
		Client: clover.Client{
			ID:     "c1",
			Scopes: []string{"s1", "s2"},
			Public: false,
		},
	}
	expAtCtx := &clover.AccessTokenContext{
		Client: clover.Client{
			ID:     "c1",
			Scopes: []string{"s1", "s2"},
			Public: false,
		},
		Scopes:              []string{"s1", "s2"},
		AccessTokenLifespan: 15,
	}

	t.Run("Success", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}

		err := g.Validate(atCtx, tokenManager)
		require.NoError(t, err)
		require.Equal(t, expAtCtx, atCtx)
		tokenManager.AssertExpectations(t)
	})

	t.Run("Client public not allowed", func(t *testing.T) {
		atCtx := *atCtx
		atCtx.Client.Public = true
		tokenManager := &mocks.TokenManager{}

		err := g.Validate(&atCtx, tokenManager)
		require.Equal(t, `The client is public and thus not allowed to use grant type client_credentials`, err.Error())
		tokenManager.AssertExpectations(t)
	})
}

func testClientGrantTypeCreateAccessToken(t *testing.T) {
	g := clover.ClientCredentialsGrantType{}

	atCtx := &clover.AccessTokenContext{
		Client:              clover.Client{ID: "c1", IncludeRefreshToken: true},
		Scopes:              []string{"s1", "s2"},
		AccessTokenLifespan: 15,
	}
	at := &clover.AccessToken{
		AccessToken: "accesstoken",
		Scopes:      []string{"s1", "s2"},
	}

	t.Run("Success", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}
		expRes := &clover.AccessTokenRes{
			AccessToken: "accesstoken",
			TokenType:   "bearer",
			ExpiresIn:   15,
			Scope:       "s1 s2",
		}

		tokenManager.On("GenerateAccessToken", atCtx, false).Return(at, nil, nil)

		res, err := g.CreateAccessToken(atCtx, tokenManager)
		require.NoError(t, err)
		require.Equal(t, expRes, res)
		tokenManager.AssertExpectations(t)
	})

	t.Run("GenerateAccessToken Failed", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}
		tokenManager.On("GenerateAccessToken", atCtx, false).Return(nil, nil, errors.New("generate failed"))

		res, err := g.CreateAccessToken(atCtx, tokenManager)
		require.Equal(t, "generate failed", err.Error())
		require.Nil(t, res)
		tokenManager.AssertExpectations(t)
	})
}
