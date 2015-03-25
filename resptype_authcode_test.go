package clover

import (
	"errors"
	"github.com/plimble/unik/mock_unik"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCodeRespType struct {
	store *Mockallstore
	unik  *mock_unik.MockGenerator
}

func setupCodeRespType() (*codeRespType, *mockCodeRespType) {
	store := NewMockallstore()
	unik := mock_unik.NewMockGenerator()

	config := NewAuthConfig(store)
	config.AuthCodeStore = store

	mock := &mockCodeRespType{store, unik}
	rt := newCodeRespType(config, unik)
	return rt, mock
}

func TestCodeRespType_GetAuthResponse(t *testing.T) {
	rt, m := setupCodeRespType()
	rt.config.AuthCodeLifetime = 60

	ar := &authorizeRequest{
		redirectURI:  "http://localhost/redirect",
		responseType: RESP_TYPE_CODE,
		state:        "123",
	}

	client := &DefaultClient{
		ClientID: "123",
		UserID:   "abc",
	}

	scopes := []string{"1", "2"}

	m.unik.On("Generate").Return("xyz")
	expAc := &AuthorizeCode{
		Code:        rt.generateAuthCode(),
		ClientID:    client.ClientID,
		UserID:      client.UserID,
		Expires:     addSecondUnix(60),
		Scope:       scopes,
		RedirectURI: ar.redirectURI,
	}

	m.unik.On("Generate").Return("xyz")
	m.store.On("SetAuthorizeCode", expAc).Return(nil)

	resp := rt.GetAuthResponse(ar, client, scopes)
	assert.False(t, resp.IsError())
	assert.Equal(t, 302, resp.code)
	assert.Equal(t, "http://localhost/redirect", resp.redirectURI)
	assert.Equal(t, map[string]interface{}{"code": expAc.Code, "state": "123"}, resp.data)
	assert.False(t, resp.isFragment)
}

func TestCodeRespType_CreateAuthCode(t *testing.T) {
	rt, m := setupCodeRespType()
	rt.config.AuthCodeLifetime = 60

	client := &DefaultClient{
		ClientID: "123",
		UserID:   "abc",
	}

	m.unik.On("Generate").Return("xyz")
	expAc := &AuthorizeCode{
		Code:        rt.generateAuthCode(),
		ClientID:    client.ClientID,
		UserID:      client.UserID,
		Expires:     addSecondUnix(60),
		Scope:       []string{"1", "2"},
		RedirectURI: "http://localhost/redirect",
	}

	m.unik.On("Generate").Return("xyz")
	m.store.On("SetAuthorizeCode", expAc).Return(nil)

	ac, resp := rt.createAuthCode(client, []string{"1", "2"}, "http://localhost/redirect")
	assert.Equal(t, expAc, ac)
	assert.Nil(t, resp)
}

func TestCodeRespType_CreateAuthCode_WithError(t *testing.T) {
	rt, m := setupCodeRespType()
	rt.config.AuthCodeLifetime = 60

	client := &DefaultClient{
		ClientID: "123",
		UserID:   "abc",
	}

	m.unik.On("Generate").Return("xyz")
	expAc := &AuthorizeCode{
		Code:        rt.generateAuthCode(),
		ClientID:    client.ClientID,
		UserID:      client.UserID,
		Expires:     addSecondUnix(60),
		Scope:       []string{"1", "2"},
		RedirectURI: "http://localhost/redirect",
	}

	m.unik.On("Generate").Return("xyz")
	m.store.On("SetAuthorizeCode", expAc).Return(errors.New("error"))

	ac, resp := rt.createAuthCode(client, []string{"1", "2"}, "http://localhost/redirect")
	assert.Nil(t, ac)
	assert.Equal(t, errInternal("error"), resp)
}

func TestCodeRespType_CreateRespData(t *testing.T) {
	rt, _ := setupCodeRespType()

	testCases := []struct {
		code     string
		state    string
		respData respData
	}{
		//with out state
		{code: "123", state: "", respData: respData{"code": "123"}},
		//with state
		{code: "123", state: "123", respData: respData{"code": "123", "state": "123"}},
	}

	for _, testCase := range testCases {
		data := rt.createRespData(testCase.code, testCase.state)
		assert.Equal(t, testCase.respData, data)
	}
}
