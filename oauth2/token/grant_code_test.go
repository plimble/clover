package token_test

// import (
// 	"errors"
// 	"testing"

// 	"github.com/plimble/clover"
// 	"github.com/plimble/clover/mocks"
// 	"github.com/stretchr/testify/require"
// )

// func TestAuthorizeCodeGrantType(t *testing.T) {
// 	t.Run("Name", testAuthorizeCodeGrantTypeName)
// 	t.Run("Validate", testAuthorizeCodeGrantTypeValidate)
// 	t.Run("CreateAccessToken", testAuthorizeCodeGrantTypeCreateAccessToken)
// }

// func testAuthorizeCodeGrantTypeName(t *testing.T) {
// 	g := clover.AuthorizeCodeGrantType{}
// 	require.Equal(t, "authorization_code", g.Name())
// }

// func testAuthorizeCodeGrantTypeValidate(t *testing.T) {
// 	g := clover.AuthorizeCodeGrantType{
// 		AccessTokenLifespan:  15,
// 		RefreshTokenLifespan: 16,
// 	}

// 	atCtx := &clover.AccessTokenContext{
// 		Client:      clover.Client{ID: "c1"},
// 		Code:        "code1234",
// 		RedirectURI: "http://example.com",
// 	}
// 	expAtCtx := &clover.AccessTokenContext{
// 		Client:               clover.Client{ID: "c1"},
// 		Code:                 "code1234",
// 		UserID:               "u1",
// 		Scopes:               []string{"s1", "s2"},
// 		AccessTokenLifespan:  15,
// 		RefreshTokenLifespan: 16,
// 		RedirectURI:          "http://example.com",
// 	}
// 	authCode := &clover.AuthorizeCode{
// 		ClientID:    "c1",
// 		Code:        "code1234",
// 		UserID:      "u1",
// 		Scopes:      []string{"s1", "s2"},
// 		Expired:     addSecondUnix(15),
// 		RedirectURI: "http://example.com",
// 	}

// 	t.Run("Success", func(t *testing.T) {
// 		tokenManager := &mocks.TokenManager{}

// 		tokenManager.On("GetAuthorizeCode", "code1234").Return(authCode, nil)

// 		err := g.Validate(atCtx, tokenManager)
// 		require.NoError(t, err)
// 		require.Equal(t, expAtCtx, atCtx)
// 		tokenManager.AssertExpectations(t)
// 	})

// 	t.Run("AuthorizeCode Required", func(t *testing.T) {
// 		atCtx := *atCtx
// 		atCtx.Code = ""
// 		tokenManager := &mocks.TokenManager{}

// 		err := g.Validate(&atCtx, tokenManager)
// 		require.Equal(t, `Missing parameter: "code" is required`, err.Error())
// 		tokenManager.AssertExpectations(t)
// 	})

// 	t.Run("RedirectURI Required", func(t *testing.T) {
// 		atCtx := *atCtx
// 		atCtx.RedirectURI = ""
// 		tokenManager := &mocks.TokenManager{}

// 		err := g.Validate(&atCtx, tokenManager)
// 		require.Equal(t, `Missing parameter: "redirect_uri" is required`, err.Error())
// 		tokenManager.AssertExpectations(t)
// 	})

// 	t.Run("Client mismatched", func(t *testing.T) {
// 		atCtx := *atCtx
// 		atCtx.Client.ID = "c2"
// 		tokenManager := &mocks.TokenManager{}

// 		tokenManager.On("GetAuthorizeCode", "code1234").Return(authCode, nil)

// 		err := g.Validate(&atCtx, tokenManager)
// 		require.Equal(t, `Client mismatched`, err.Error())
// 		tokenManager.AssertExpectations(t)
// 	})

// 	t.Run("RedirectURI mismatched", func(t *testing.T) {
// 		atCtx := *atCtx
// 		atCtx.RedirectURI = "http://abc.com"
// 		tokenManager := &mocks.TokenManager{}

// 		tokenManager.On("GetAuthorizeCode", "code1234").Return(authCode, nil)

// 		err := g.Validate(&atCtx, tokenManager)
// 		require.Equal(t, `The redirect URI is missing or do not match`, err.Error())
// 		tokenManager.AssertExpectations(t)
// 	})

// 	t.Run("Code expired", func(t *testing.T) {
// 		authCode := *authCode
// 		authCode.Expired = addSecondUnix(-100)

// 		tokenManager := &mocks.TokenManager{}

// 		tokenManager.On("GetAuthorizeCode", "code1234").Return(&authCode, nil)

// 		err := g.Validate(atCtx, tokenManager)
// 		require.Equal(t, `The authorization code has expire`, err.Error())
// 		tokenManager.AssertExpectations(t)
// 	})
// }

// func testAuthorizeCodeGrantTypeCreateAccessToken(t *testing.T) {
// 	g := clover.AuthorizeCodeGrantType{}

// 	atCtx := &clover.AccessTokenContext{
// 		Client:               clover.Client{ID: "c1", IncludeRefreshToken: true},
// 		UserID:               "u1",
// 		Scopes:               []string{"s1", "s2"},
// 		AccessTokenLifespan:  15,
// 		RefreshTokenLifespan: 16,
// 	}
// 	at := &clover.AccessToken{
// 		AccessToken: "accesstoken",
// 		Scopes:      []string{"s1", "s2"},
// 		UserID:      "u1",
// 	}
// 	rf := &clover.RefreshToken{
// 		RefreshToken: "refreshtoken",
// 	}

// 	t.Run("Success With RefreshToken", func(t *testing.T) {
// 		tokenManager := &mocks.TokenManager{}
// 		expRes := &clover.AccessTokenRes{
// 			AccessToken:  "accesstoken",
// 			TokenType:    "bearer",
// 			ExpiresIn:    15,
// 			Scope:        "s1 s2",
// 			UserID:       "u1",
// 			RefreshToken: "refreshtoken",
// 		}

// 		tokenManager.On("GenerateAccessToken", atCtx, true).Return(at, rf, nil)

// 		res, err := g.CreateAccessToken(atCtx, tokenManager)
// 		require.NoError(t, err)
// 		require.Equal(t, expRes, res)
// 		tokenManager.AssertExpectations(t)
// 	})

// 	t.Run("Success Without RefreshToken", func(t *testing.T) {
// 		atCtx := *atCtx
// 		atCtx.Client.IncludeRefreshToken = false

// 		tokenManager := &mocks.TokenManager{}
// 		expRes := &clover.AccessTokenRes{
// 			AccessToken: "accesstoken",
// 			TokenType:   "bearer",
// 			ExpiresIn:   15,
// 			Scope:       "s1 s2",
// 			UserID:      "u1",
// 		}

// 		tokenManager.On("GenerateAccessToken", &atCtx, false).Return(at, rf, nil)

// 		res, err := g.CreateAccessToken(&atCtx, tokenManager)
// 		require.NoError(t, err)
// 		require.Equal(t, expRes, res)
// 		tokenManager.AssertExpectations(t)
// 	})

// 	t.Run("GenerateAccessToken Failed", func(t *testing.T) {
// 		tokenManager := &mocks.TokenManager{}
// 		tokenManager.On("GenerateAccessToken", atCtx, true).Return(nil, nil, errors.New("generate failed"))

// 		res, err := g.CreateAccessToken(atCtx, tokenManager)
// 		require.Equal(t, "generate failed", err.Error())
// 		require.Nil(t, res)
// 		tokenManager.AssertExpectations(t)
// 	})
// }
