package token_test

import (
	"errors"
	"testing"

	"github.com/plimble/clover"
	"github.com/plimble/clover/mocks"
	"github.com/stretchr/testify/require"
)

func TestRefreshGrantType(t *testing.T) {
	t.Run("Name", testRefreshGrantTypeName)
	t.Run("Validate", testRefreshGrantTypeValidate)
	t.Run("CreateAccessToken", testRefreshGrantTypeCreateAccessToken)
}

func testRefreshGrantTypeName(t *testing.T) {
	g := clover.RefreshTokenGrantType{}
	require.Equal(t, "refresh_token", g.Name())
}

func testRefreshGrantTypeValidate(t *testing.T) {
	g := clover.RefreshTokenGrantType{}

	atCtx := &clover.AccessTokenContext{
		RefreshToken: "token",
		Client:       clover.Client{ID: "c1"},
	}
	expAtCtx := &clover.AccessTokenContext{
		RefreshToken:         "token",
		Client:               clover.Client{ID: "c1"},
		UserID:               "u1",
		Scopes:               []string{"s1", "s2"},
		AccessTokenLifespan:  15,
		RefreshTokenLifespan: 16,
	}
	rf := &clover.RefreshToken{
		RefreshToken:         "token",
		ClientID:             "c1",
		UserID:               "u1",
		Scopes:               []string{"s1", "s2"},
		AccessTokenLifespan:  15,
		RefreshTokenLifespan: 16,
		Expired:              addSecondUnix(16),
	}

	t.Run("Success", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}

		tokenManager.On("GetRefreshToken", "token").Return(rf, nil)

		err := g.Validate(atCtx, tokenManager)
		require.NoError(t, err)
		require.Equal(t, expAtCtx, atCtx)
		tokenManager.AssertExpectations(t)
	})

	t.Run("RefreshToken Required", func(t *testing.T) {
		atCtx := *atCtx
		atCtx.RefreshToken = ""
		tokenManager := &mocks.TokenManager{}

		err := g.Validate(&atCtx, tokenManager)
		require.Equal(t, `Missing parameter: "refresh_token" is required`, err.Error())
		tokenManager.AssertExpectations(t)
	})

	t.Run("Invalid RefreshToken", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}

		tokenManager.On("GetRefreshToken", "token").Return(nil, errors.New("invalid token"))

		err := g.Validate(atCtx, tokenManager)
		require.Equal(t, `The refresh token provided is invalid`, err.Error())
		tokenManager.AssertExpectations(t)
	})

	t.Run("Invalid Client MisMatched", func(t *testing.T) {
		rf := *rf
		rf.ClientID = "c2"
		tokenManager := &mocks.TokenManager{}

		tokenManager.On("GetRefreshToken", "token").Return(&rf, nil)

		err := g.Validate(atCtx, tokenManager)
		require.Equal(t, `Client mismatched`, err.Error())
		tokenManager.AssertExpectations(t)
	})

	t.Run("RefreshToken Expired", func(t *testing.T) {
		rf := *rf
		rf.Expired = addSecondUnix(-100)
		tokenManager := &mocks.TokenManager{}

		tokenManager.On("GetRefreshToken", "token").Return(&rf, nil)

		err := g.Validate(atCtx, tokenManager)
		require.Equal(t, `Refresh token has expired`, err.Error())
		tokenManager.AssertExpectations(t)
	})
}

func testRefreshGrantTypeCreateAccessToken(t *testing.T) {
	g := clover.RefreshTokenGrantType{}

	atCtx := &clover.AccessTokenContext{
		Client:               clover.Client{ID: "c1", IncludeRefreshToken: true},
		UserID:               "u1",
		Scopes:               []string{"s1", "s2"},
		AccessTokenLifespan:  15,
		RefreshTokenLifespan: 16,
		RefreshToken:         "old_refreshtoken",
	}
	at := &clover.AccessToken{
		AccessToken: "accesstoken",
		Scopes:      []string{"s1", "s2"},
		UserID:      "u1",
	}
	rf := &clover.RefreshToken{
		RefreshToken: "new_refreshtoken",
	}

	t.Run("Success", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}
		expRes := &clover.AccessTokenRes{
			AccessToken:  "accesstoken",
			TokenType:    "bearer",
			ExpiresIn:    15,
			Scope:        "s1 s2",
			UserID:       "u1",
			RefreshToken: "new_refreshtoken",
		}

		tokenManager.On("GenerateAccessToken", atCtx, true).Return(at, rf, nil)
		tokenManager.On("DeleteRefreshToken", "old_refreshtoken").Return(nil)

		res, err := g.CreateAccessToken(atCtx, tokenManager)
		require.NoError(t, err)
		require.Equal(t, expRes, res)
		tokenManager.AssertExpectations(t)
	})

	t.Run("GenerateAccessToken Failed", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}
		tokenManager.On("GenerateAccessToken", atCtx, true).Return(nil, nil, errors.New("generate failed"))

		res, err := g.CreateAccessToken(atCtx, tokenManager)
		require.Equal(t, "generate failed", err.Error())
		require.Nil(t, res)
		tokenManager.AssertExpectations(t)
	})

	t.Run("DeleteRefreshToken Failed", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}
		tokenManager.On("GenerateAccessToken", atCtx, true).Return(at, rf, nil)
		tokenManager.On("DeleteRefreshToken", "old_refreshtoken").Return(errors.New("delete failed"))

		res, err := g.CreateAccessToken(atCtx, tokenManager)
		require.Equal(t, "delete failed", err.Error())
		require.Nil(t, res)
		tokenManager.AssertExpectations(t)
	})
}
