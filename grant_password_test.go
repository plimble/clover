package clover_test

import (
	"errors"
	"testing"

	"github.com/plimble/clover"
	"github.com/plimble/clover/mocks"
	"github.com/stretchr/testify/require"
)

func TestPasswordGrantType(t *testing.T) {
	t.Run("Name", testPasswordGrantTypeName)
	t.Run("Validate", testPasswordGrantTypeValidate)
	t.Run("CreateAccessToken", testPasswordGrantTypeCreateAccessToken)
}

func testPasswordGrantTypeName(t *testing.T) {
	g := clover.PassowrdGrantType{}
	require.Equal(t, "password", g.Name())
}

func testPasswordGrantTypeValidate(t *testing.T) {
	g := clover.PassowrdGrantType{
		AccessTokenLifespan:  15,
		RefreshTokenLifespan: 16,
	}

	atCtx := &clover.AccessTokenContext{
		Client:   clover.Client{ID: "c1"},
		Username: "u1",
		Password: "pwd",
	}
	expAtCtx := &clover.AccessTokenContext{
		Client:               clover.Client{ID: "c1"},
		Username:             "u1",
		Password:             "pwd",
		UserID:               "u1",
		Scopes:               []string{"s1", "s2"},
		AccessTokenLifespan:  15,
		RefreshTokenLifespan: 16,
	}
	u := &clover.User{
		ID:     "u1",
		Scopes: []string{"s1", "s2"},
	}

	t.Run("Success", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}
		userService := &mocks.UserService{}
		g.UserService = userService

		userService.On("GetUser", "u1", "pwd").Return(u, nil)

		err := g.Validate(atCtx, tokenManager)
		require.NoError(t, err)
		require.Equal(t, expAtCtx, atCtx)
		tokenManager.AssertExpectations(t)
	})

	t.Run("Username password Required", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}
		userService := &mocks.UserService{}
		g.UserService = userService
		atCtx := *atCtx
		atCtx.Username = ""
		atCtx.Password = ""

		err := g.Validate(&atCtx, tokenManager)
		require.Equal(t, `Missing parameters: "username" and "password" required`, err.Error())
		tokenManager.AssertExpectations(t)
	})

	t.Run("Username password Required", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}
		userService := &mocks.UserService{}
		g.UserService = userService

		userService.On("GetUser", "u1", "pwd").Return(nil, errors.New("invalid user"))

		err := g.Validate(atCtx, tokenManager)
		require.Equal(t, `Invalid username or password`, err.Error())
		tokenManager.AssertExpectations(t)
	})
}

func testPasswordGrantTypeCreateAccessToken(t *testing.T) {
	g := clover.PassowrdGrantType{}

	atCtx := &clover.AccessTokenContext{
		Client:               clover.Client{ID: "c1", IncludeRefreshToken: true},
		UserID:               "u1",
		Scopes:               []string{"s1", "s2"},
		AccessTokenLifespan:  15,
		RefreshTokenLifespan: 16,
	}
	at := &clover.AccessToken{
		AccessToken: "accesstoken",
		Scopes:      []string{"s1", "s2"},
		UserID:      "u1",
	}
	rf := &clover.RefreshToken{
		RefreshToken: "refreshtoken",
	}

	t.Run("Success With RefreshToken", func(t *testing.T) {
		tokenManager := &mocks.TokenManager{}
		expRes := &clover.AccessTokenRes{
			AccessToken:  "accesstoken",
			TokenType:    "bearer",
			ExpiresIn:    15,
			Scope:        "s1 s2",
			UserID:       "u1",
			RefreshToken: "refreshtoken",
		}

		tokenManager.On("GenerateAccessToken", atCtx, true).Return(at, rf, nil)

		res, err := g.CreateAccessToken(atCtx, tokenManager)
		require.NoError(t, err)
		require.Equal(t, expRes, res)
		tokenManager.AssertExpectations(t)
	})

	t.Run("Success Without RefreshToken", func(t *testing.T) {
		atCtx := *atCtx
		atCtx.Client.IncludeRefreshToken = false

		tokenManager := &mocks.TokenManager{}
		expRes := &clover.AccessTokenRes{
			AccessToken: "accesstoken",
			TokenType:   "bearer",
			ExpiresIn:   15,
			Scope:       "s1 s2",
			UserID:      "u1",
		}

		tokenManager.On("GenerateAccessToken", &atCtx, false).Return(at, rf, nil)

		res, err := g.CreateAccessToken(&atCtx, tokenManager)
		require.NoError(t, err)
		require.Equal(t, expRes, res)
		tokenManager.AssertExpectations(t)
	})

	t.Run("GenerateAccessToken Failed", func(t *testing.T) {
		atCtx := *atCtx

		tokenManager := &mocks.TokenManager{}
		tokenManager.On("GenerateAccessToken", &atCtx, true).Return(nil, nil, errors.New("generate failed"))

		res, err := g.CreateAccessToken(&atCtx, tokenManager)
		require.Equal(t, "generate failed", err.Error())
		require.Nil(t, res)
		tokenManager.AssertExpectations(t)
	})
}
