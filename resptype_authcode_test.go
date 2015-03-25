package clover

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"github.com/plimble/unik/mock_unik"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockResponseType struct {
	store *Mockallstore
}

func setUpCodeResponseType() (*codeRespType, *mockResponseType) {
	store := NewMockallstore()
	mock := &mockResponseType{store}

	config := NewAuthConfig(store)
	config.AddAuthCodeGrant(store)

	mu := mock_unik.NewMockGenerator()
	mu.On("Generate").Return("1")
	rt := newCodeRespType(config, mu)
	return rt, mock
}

func generateAuthRequest(resType string) *authorizeRequest {
	ar := &authorizeRequest{
		state:        "0",
		redirectURI:  "http://localhost",
		responseType: resType,
		scope:        "email",
		clientID:     "1001",
		client: &DefaultClient{
			ClientID:     "1",
			RedirectURI:  "http://localhost",
			ClientSecret: "xyz",
			GrantType:    []string{AUTHORIZATION_CODE},
			UserID:       "1",
			Scope:        []string{"email"},
		},
	}

	return ar
}

func genTestAuthCode(rt *codeRespType, ar *authorizeRequest) *AuthorizeCode {
	return &AuthorizeCode{
		Code:        hashCode("1"),
		ClientID:    ar.client.GetClientID(),
		UserID:      ar.client.GetUserID(),
		Expires:     addSecondUnix(rt.config.AuthCodeLifetime),
		Scope:       []string{"email"},
		RedirectURI: ar.redirectURI,
	}
}

func hashCode(code string) string {
	hasher := sha512.New()
	hasher.Write([]byte(code))
	return hex.EncodeToString(hasher.Sum(nil))
}

func TestGetAuthResponse(t *testing.T) {
	rt, mock := setUpCodeResponseType()
	ar := generateAuthRequest("code")
	ac := genTestAuthCode(rt, ar)

	mock.store.On("SetAuthorizeCode", ac).Return(nil)
	resp := rt.GetAuthResponse(ar, ar.client, ar.client.GetScope())

	assert.Equal(t, 302, resp.code)
	assert.False(t, resp.IsError())
}

func TestGetAuthResponseError(t *testing.T) {
	rt, mock := setUpCodeResponseType()
	ar := generateAuthRequest("code")
	ac := genTestAuthCode(rt, ar)

	mock.store.On("SetAuthorizeCode", ac).Return(errors.New("test"))
	resp := rt.GetAuthResponse(ar, ar.client, ar.client.GetScope())

	assert.Equal(t, 500, resp.code)
	assert.True(t, resp.IsError())
}

func TestCreateAuthCode(t *testing.T) {
	rt, mock := setUpCodeResponseType()
	ar := generateAuthRequest("code")
	ac := genTestAuthCode(rt, ar)

	mock.store.On("SetAuthorizeCode", ac).Return(nil)
	code, resp := rt.createAuthCode(ar.client, ar.client.GetScope(), ar.client.GetRedirectURI())

	assert.Nil(t, resp)
	assert.Equal(t, ac.Code, code.Code)
}

func TestCreateAuthCodeError(t *testing.T) {
	rt, mock := setUpCodeResponseType()
	ar := generateAuthRequest("code")
	ac := genTestAuthCode(rt, ar)

	mock.store.On("SetAuthorizeCode", ac).Return(errors.New("test"))
	code, resp := rt.createAuthCode(ar.client, ar.client.GetScope(), ar.client.GetRedirectURI())

	assert.Nil(t, code)
	assert.True(t, resp.IsError())
}

func TestGenerateAuthCode(t *testing.T) {
	rt, _ := setUpCodeResponseType()
	code := rt.generateAuthCode()
	assert.NotEmpty(t, code)
}
