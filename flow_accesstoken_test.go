package clover_test

import (
	"errors"
	"testing"

	"github.com/plimble/clover"
	"github.com/plimble/clover/mocks"
	"github.com/stretchr/testify/require"
)

type accessTokenFlow struct {
	tokenManager   *mocks.TokenManager
	grantTypes     map[string]clover.GrantType
	scopeValidator *mocks.ScopeValidator
	clientStorage  *mocks.ClientStorage
}

func setUpAccessTokenFlow() *accessTokenFlow {
	return &accessTokenFlow{
		tokenManager: &mocks.TokenManager{},
		grantTypes: map[string]clover.GrantType{
			"client_credentials": &clover.ClientCredentialsGrantType{
				AccessTokenLifespan: 15,
			},
			"refresh_token": &clover.RefreshTokenGrantType{},
		},
		scopeValidator: &mocks.ScopeValidator{},
		clientStorage:  &mocks.ClientStorage{},
	}
}

func TestAccessTokenFlow(t *testing.T) {
	c := &clover.Client{
		ID:                  "c1",
		Scopes:              []string{"s1", "s2"},
		Secret:              "secret",
		IncludeRefreshToken: false,
	}
	atCtx := &clover.AccessTokenContext{
		HTTPContext: clover.HTTPContext{
			AuthorizationHeader: "Basic YzE6c2VjcmV0",
			Method:              "POST",
		},
		GrantType: "client_credentials",
	}
	at := &clover.AccessToken{
		AccessToken: "accesstoken",
		Scopes:      []string{"s1", "s2"},
	}
	expRes := &clover.AccessTokenRes{
		AccessToken: "accesstoken",
		TokenType:   "bearer",
		ExpiresIn:   15,
		Scope:       "s1 s2",
	}

	t.Run("Success with client_credentials grant type", func(t *testing.T) {
		m := setUpAccessTokenFlow()
		f := clover.NewAccessTokenFlow(m.tokenManager, m.grantTypes, m.scopeValidator, m.clientStorage)

		m.clientStorage.On("GetClientWithSecret", "c1", "secret").Return(c, nil)
		m.tokenManager.On("GenerateAccessToken", atCtx, false).Return(at, nil, nil)

		res, err := f.Run(atCtx)
		require.NoError(t, err)
		require.Equal(t, expRes, res)
	})

	t.Run("Success with refreshtoken grant type", func(t *testing.T) {
		atCtx := &clover.AccessTokenContext{
			HTTPContext: clover.HTTPContext{
				AuthorizationHeader: "Basic YzE6c2VjcmV0",
				Method:              "POST",
			},
			GrantType:    "refresh_token",
			RefreshToken: "old_refreshtoken",
		}
		rt := &clover.RefreshToken{
			ClientID:            "c1",
			Expired:             addSecondUnix(15),
			Scopes:              []string{"s1", "s2"},
			AccessTokenLifespan: 15,
		}

		m := setUpAccessTokenFlow()
		f := clover.NewAccessTokenFlow(m.tokenManager, m.grantTypes, m.scopeValidator, m.clientStorage)

		m.clientStorage.On("GetClientWithSecret", "c1", "secret").Return(c, nil)
		m.scopeValidator.On("Validate", []string{"s1", "s2"}, []string{"s1", "s2"}).Return([]string{"s1", "s2"}, nil)
		m.tokenManager.On("GetRefreshToken", "old_refreshtoken").Return(rt, nil)
		m.tokenManager.On("GenerateAccessToken", atCtx, false).Return(at, nil, nil)
		m.tokenManager.On("DeleteRefreshToken", "old_refreshtoken").Return(nil)

		res, err := f.Run(atCtx)
		require.NoError(t, err)
		require.Equal(t, expRes, res)
	})

	t.Run("Error scope not supported", func(t *testing.T) {
		atCtx := &clover.AccessTokenContext{
			HTTPContext: clover.HTTPContext{
				AuthorizationHeader: "Basic YzE6c2VjcmV0",
				Method:              "POST",
			},
			GrantType:    "refresh_token",
			RefreshToken: "old_refreshtoken",
		}
		rt := &clover.RefreshToken{
			ClientID:            "c1",
			Expired:             addSecondUnix(15),
			Scopes:              []string{"s1", "s2"},
			AccessTokenLifespan: 15,
		}

		m := setUpAccessTokenFlow()
		f := clover.NewAccessTokenFlow(m.tokenManager, m.grantTypes, m.scopeValidator, m.clientStorage)

		m.clientStorage.On("GetClientWithSecret", "c1", "secret").Return(c, nil)
		m.scopeValidator.On("Validate", []string{"s1", "s2"}, []string{"s1", "s2"}).Return(nil, errors.New("not supported"))
		m.tokenManager.On("GetRefreshToken", "old_refreshtoken").Return(rt, nil)

		res, err := f.Run(atCtx)
		require.Equal(t, `The scope requested is invalid for this request`, err.Error())
		require.Nil(t, res)
	})

	t.Run("Method not support", func(t *testing.T) {
		atCtx := *atCtx
		atCtx.Method = "GET"
		m := setUpAccessTokenFlow()
		f := clover.NewAccessTokenFlow(m.tokenManager, m.grantTypes, m.scopeValidator, m.clientStorage)

		res, err := f.Run(&atCtx)
		require.Equal(t, "The request method must be POST when requesting an access token", err.Error())
		require.Nil(t, res)
	})

	t.Run("Granttype required", func(t *testing.T) {
		atCtx := *atCtx
		atCtx.GrantType = ""
		m := setUpAccessTokenFlow()
		f := clover.NewAccessTokenFlow(m.tokenManager, m.grantTypes, m.scopeValidator, m.clientStorage)

		res, err := f.Run(&atCtx)
		require.Equal(t, `Missing parameters: "grant_type" required`, err.Error())
		require.Nil(t, res)
	})

	t.Run("Granttype not support", func(t *testing.T) {
		atCtx := *atCtx
		atCtx.GrantType = "abc"
		m := setUpAccessTokenFlow()
		f := clover.NewAccessTokenFlow(m.tokenManager, m.grantTypes, m.scopeValidator, m.clientStorage)

		res, err := f.Run(&atCtx)
		require.Equal(t, `Grant type "abc" not supported`, err.Error())
		require.Nil(t, res)
	})

	t.Run("Error get client", func(t *testing.T) {
		m := setUpAccessTokenFlow()
		f := clover.NewAccessTokenFlow(m.tokenManager, m.grantTypes, m.scopeValidator, m.clientStorage)

		m.clientStorage.On("GetClientWithSecret", "c1", "secret").Return(nil, errors.New("error get client"))

		res, err := f.Run(atCtx)
		require.Equal(t, `Invalid client`, err.Error())
		require.Nil(t, res)
	})

	t.Run("Error validate grant", func(t *testing.T) {
		c := *c
		c.Public = true
		m := setUpAccessTokenFlow()
		f := clover.NewAccessTokenFlow(m.tokenManager, m.grantTypes, m.scopeValidator, m.clientStorage)

		m.clientStorage.On("GetClientWithSecret", "c1", "secret").Return(&c, nil)

		res, err := f.Run(atCtx)
		require.Equal(t, `The client is public and thus not allowed to use grant type client_credentials`, err.Error())
		require.Nil(t, res)
	})

	t.Run("Error generate token", func(t *testing.T) {
		m := setUpAccessTokenFlow()
		f := clover.NewAccessTokenFlow(m.tokenManager, m.grantTypes, m.scopeValidator, m.clientStorage)

		m.clientStorage.On("GetClientWithSecret", "c1", "secret").Return(c, nil)
		m.tokenManager.On("GenerateAccessToken", atCtx, false).Return(nil, nil, errors.New("error gen"))

		res, err := f.Run(atCtx)
		require.Equal(t, `Internal Server Error`, err.Error())
		require.Nil(t, res)
	})
}
