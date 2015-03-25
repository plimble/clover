package clover

import (
	"errors"
	"github.com/plimble/unik/mock_unik"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockTokenRespType struct {
	store *Mockallstore
	unik  *mock_unik.MockGenerator
}

func setupTokenResponseType() (*tokenRespType, *mockTokenRespType) {
	store := NewMockallstore()
	unik := mock_unik.NewMockGenerator()

	config := NewAuthConfig(store)
	config.RefreshTokenStore = store
	config.AuthCodeStore = store

	mock := &mockTokenRespType{store, unik}
	rt := newTokenRespType(config, unik)
	return rt, mock
}

func TestTokenRespType_GetAuthResponse(t *testing.T) {
	rt, mock := setupTokenResponseType()

	ar := &authorizeRequest{
		redirectURI:  "http://localhost/redirect",
		responseType: RESP_TYPE_TOKEN,
		state:        "123",
	}

	at := &AccessToken{
		ClientID:    "1001",
		UserID:      "1",
		AccessToken: "123",
		Scope:       []string{"email"},
		Expires:     addSecondUnix(rt.config.AccessLifeTime),
	}

	client := &DefaultClient{
		ClientID: "1001",
		UserID:   "1",
	}

	mock.unik.On("Generate").Return("123")
	mock.store.On("SetAccessToken", at).Return(nil)
	resp := rt.GetAuthResponse(ar, client, []string{"email"})
	assert.False(t, resp.IsError())
	assert.Equal(t, 302, resp.code)
	assert.Equal(t, at.AccessToken, resp.data["access_token"])
}

func TestTokenRespType_GetAuthResponse_WithError(t *testing.T) {
	rt, mock := setupTokenResponseType()

	ar := &authorizeRequest{
		redirectURI:  "http://localhost/redirect",
		responseType: RESP_TYPE_TOKEN,
		state:        "123",
	}

	at := &AccessToken{
		ClientID:    "1001",
		UserID:      "1",
		AccessToken: "123",
		Scope:       []string{"email"},
		Expires:     addSecondUnix(rt.config.AccessLifeTime),
	}

	client := &DefaultClient{
		ClientID: "1001",
		UserID:   "1",
	}

	mock.unik.On("Generate").Return("123")
	mock.store.On("SetAccessToken", at).Return(errors.New("test"))
	resp := rt.GetAuthResponse(ar, client, []string{"email"})
	assert.True(t, resp.IsError())
	assert.Equal(t, 500, resp.code)
}

func TestTokenRespType_GetAccessToken(t *testing.T) {
	rt, mock := setupTokenResponseType()
	expTime := addSecondUnix(rt.config.AccessLifeTime)
	at := &AccessToken{
		AccessToken: "123",
		Expires:     expTime,
	}

	r := &RefreshToken{
		RefreshToken: "123",
		Expires:      expTime,
	}

	td := &TokenData{
		GrantData: &GrantData{
			RefreshToken: "123",
		}}

	mock.unik.On("Generate").Return("123")
	mock.store.On("SetAccessToken", at).Return(nil)
	mock.store.On("SetRefreshToken", r).Return(nil)

	resp := rt.GetAccessToken(td, true)
	assert.False(t, resp.IsError())
	assert.Equal(t, 200, resp.code)
	assert.Equal(t, "123", resp.data["access_token"])
}

func TestTokenRespType_GetAccessToken_WithError_SetToken(t *testing.T) {
	rt, mock := setupTokenResponseType()
	expTime := addSecondUnix(rt.config.AccessLifeTime)
	at := &AccessToken{
		AccessToken: "123",
		Expires:     expTime,
	}

	r := &RefreshToken{
		RefreshToken: "123",
		Expires:      expTime,
	}

	td := &TokenData{
		GrantData: &GrantData{
			RefreshToken: "123",
		}}

	mock.unik.On("Generate").Return("123")
	mock.store.On("SetAccessToken", at).Return(errors.New("test"))
	mock.store.On("SetRefreshToken", r).Return(nil)

	resp := rt.GetAccessToken(td, true)
	assert.True(t, resp.IsError())
	assert.Equal(t, 500, resp.code)
}

func TestTokenRespType_GetAccessToken_WithError_RefreshToken(t *testing.T) {
	rt, mock := setupTokenResponseType()
	expTime := addSecondUnix(rt.config.AccessLifeTime)
	at := &AccessToken{
		AccessToken: "123",
		Expires:     expTime,
	}

	r := &RefreshToken{
		RefreshToken: "123",
		Expires:      expTime,
	}

	td := &TokenData{
		GrantData: &GrantData{
			RefreshToken: "123",
		}}

	mock.unik.On("Generate").Return("123")
	mock.store.On("SetAccessToken", at).Return(nil)
	mock.store.On("SetRefreshToken", r).Return(errors.New("test"))

	resp := rt.GetAccessToken(td, true)
	assert.True(t, resp.IsError())
	assert.Equal(t, 500, resp.code)
}

func TestTokenRespType_CreateToken(t *testing.T) {
	rt, mock := setupTokenResponseType()
	expAt := &AccessToken{
		AccessToken: "123",
		ClientID:    "1001",
		UserID:      "1",
		Expires:     addSecondUnix(rt.config.AccessLifeTime),
		Scope:       []string{"email"},
	}

	mock.unik.On("Generate").Return(expAt.AccessToken)
	mock.store.On("SetAccessToken", expAt).Return(nil)
	at, resp := rt.createToken(expAt.ClientID, expAt.UserID, expAt.Scope)
	assert.Nil(t, resp)
	assert.Equal(t, expAt, at)
}

func TestTokenRespType_CreateToken_WithError(t *testing.T) {
	rt, mock := setupTokenResponseType()
	expAt := &AccessToken{
		AccessToken: "123",
		ClientID:    "1001",
		UserID:      "1",
		Expires:     addSecondUnix(rt.config.AccessLifeTime),
		Scope:       []string{"email"},
	}

	mock.unik.On("Generate").Return(expAt.AccessToken)
	mock.store.On("SetAccessToken", expAt).Return(errors.New("test"))
	at, resp := rt.createToken(expAt.ClientID, expAt.UserID, expAt.Scope)
	assert.True(t, resp.IsError())
	assert.Equal(t, 500, resp.code)
	assert.Nil(t, at)
}

func TestTokenRespType_CreateRefreshToken(t *testing.T) {
	rt, mock := setupTokenResponseType()
	at := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(-1),
		ClientID:    "1",
		UserID:      "1",
		Scope:       []string{"emails"},
	}

	r := &RefreshToken{
		RefreshToken: "123",
		ClientID:     at.ClientID,
		UserID:       at.UserID,
		Expires:      addSecondUnix(rt.config.RefreshTokenLifetime),
		Scope:        at.Scope,
	}

	mock.unik.On("Generate").Return("123")
	mock.store.On("SetRefreshToken", r).Return(nil)
	str, resp := rt.createRefreshToken(at, true)
	assert.Nil(t, resp)
	assert.Equal(t, r.RefreshToken, str)
}

func TestTokenRespType_CreateRefreshToken_WithError(t *testing.T) {
	rt, mock := setupTokenResponseType()
	at := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(-1),
		ClientID:    "1",
		UserID:      "1",
		Scope:       []string{"emails"},
	}

	r := &RefreshToken{
		RefreshToken: "123",
		ClientID:     at.ClientID,
		UserID:       at.UserID,
		Expires:      addSecondUnix(rt.config.RefreshTokenLifetime),
		Scope:        at.Scope,
	}

	mock.unik.On("Generate").Return("123")
	mock.store.On("SetRefreshToken", r).Return(errors.New("test"))
	str, resp := rt.createRefreshToken(at, true)
	assert.True(t, resp.IsError())
	assert.Equal(t, 500, resp.code)
	assert.Empty(t, str)
}

func TestTokenRespType_CreateRespData(t *testing.T) {
	rt, _ := setupTokenResponseType()
	testCases := []struct {
		token     string
		expiresIn int
		scopes    []string
		refresh   string
		state     string
		respData  respData
	}{
		{token: "123", expiresIn: 1, scopes: []string{"email"}, refresh: "1", state: "1", respData: respData{"access_token": "123", "token_type": "bearer", "expires_in": 3600, "scope": "email", "refresh_token": "1", "state": "1"}},
		{token: "123", expiresIn: 1, scopes: []string{"email"}, refresh: "", state: "1", respData: respData{"access_token": "123", "token_type": "bearer", "expires_in": 3600, "scope": "email", "state": "1"}},
		{token: "123", expiresIn: 1, scopes: []string{"email"}, refresh: "1", state: "", respData: respData{"access_token": "123", "token_type": "bearer", "expires_in": 3600, "scope": "email", "refresh_token": "1"}},
	}
	for _, testCase := range testCases {
		data := rt.createRespData(testCase.token, testCase.expiresIn, testCase.scopes, testCase.refresh, testCase.state)
		assert.Equal(t, testCase.respData, data)
	}
}
