package clover

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockTokenCtrl struct {
	store         *Mockallstore
	tokenRespType *MockResponseType
}

func setupTokenCtrl() (*tokenController, *mockTokenCtrl) {
	store := NewMockallstore()
	config := NewAuthConfig(store)
	config.AddClientGrant()
	config.AddPasswordGrant(store)
	config.AddRefreshGrant(store)
	config.AddAuthCodeGrant(store)

	tokenRespType := NewMockResponseType()
	mock := &mockTokenCtrl{store, tokenRespType}
	ctrl := newTokenController(config, tokenRespType)

	return ctrl, mock
}

func TestTokenCtrl_ValidateGrantType(t *testing.T) {
	ctrl, _ := setupTokenCtrl()

	testCases := []struct {
		tr        *TokenRequest
		grantName string
		resp      *Response
	}{
		{&TokenRequest{GrantType: CLIENT_CREDENTIALS}, CLIENT_CREDENTIALS, nil},
		{&TokenRequest{GrantType: PASSWORD}, PASSWORD, nil},
		{&TokenRequest{GrantType: REFRESH_TOKEN}, REFRESH_TOKEN, nil},
		{&TokenRequest{GrantType: AUTHORIZATION_CODE}, AUTHORIZATION_CODE, nil},
		{&TokenRequest{GrantType: "xxx"}, "xxx", errUnSupportedGrantType},
	}

	for _, testCase := range testCases {
		grant, resp := ctrl.validateGrantType(testCase.tr)
		assert.Equal(t, testCase.resp, resp)
		if resp == nil {
			assert.Equal(t, testCase.grantName, grant.GetGrantType())
		}
	}
}

func TestTokenCtrl_ValidateClient(t *testing.T) {
	ctrl, m := setupTokenCtrl()

	tr := &TokenRequest{
		ClientID: "1234",
	}

	grantData := &GrantData{
		ClientID: "1234",
	}

	client := &DefaultClient{
		ClientID: "1234",
	}

	m.store.On("GetClient", tr.ClientID).Return(client, nil)

	resp := ctrl.validateClient(tr, grantData)
	m.store.AssertExpectations(t)
	assert.Nil(t, resp)
}

func TestTokenCtrl_ValidateClient_WithNotFound(t *testing.T) {
	ctrl, m := setupTokenCtrl()

	tr := &TokenRequest{
		ClientID: "1234",
	}

	grantData := &GrantData{
		ClientID: "1234",
	}

	m.store.On("GetClient", tr.ClientID).Return(nil, errors.New("not found"))

	resp := ctrl.validateClient(tr, grantData)
	m.store.AssertExpectations(t)
	assert.Equal(t, errInvalidClientID, resp)
}

func TestTokenCtrl_ValidateClient_NotMatchWithClient(t *testing.T) {
	ctrl, m := setupTokenCtrl()

	tr := &TokenRequest{
		ClientID: "1234",
	}

	grantData := &GrantData{
		ClientID: "4321",
	}

	client := &DefaultClient{
		ClientID: "1234",
	}

	m.store.On("GetClient", tr.ClientID).Return(client, nil)

	resp := ctrl.validateClient(tr, grantData)
	m.store.AssertExpectations(t)
	assert.Equal(t, errInvalidClientCredentail, resp)
}

func TestTokenCtrl_ValidateScope(t *testing.T) {
	ctrl, _ := setupTokenCtrl()
	ctrl.config.DefaultScopes = []string{"3", "4"}

	testCases := []struct {
		tr        *TokenRequest
		grantData *GrantData
		scope     []string
		resp      *Response
	}{
		//scope > 0 grant scope > 0 and not in scope
		{&TokenRequest{Scope: "1 2"}, &GrantData{Scope: []string{"1"}}, nil, errInvalidScopeRequest},
		//scope > 0 grant scope > 0 and in scope
		{&TokenRequest{Scope: "1 2"}, &GrantData{Scope: []string{"1", "2"}}, []string{"1", "2"}, nil},
		//scope > 0 grant scope == 0
		{&TokenRequest{Scope: "1 2"}, &GrantData{Scope: nil}, nil, errUnSupportedScope},
		//scope == 0 grant scope > 0
		{&TokenRequest{}, &GrantData{Scope: []string{"1", "2"}}, []string{"1", "2"}, nil},
		//scope == 0 grant scope == 0 get defualt scope
		{&TokenRequest{}, &GrantData{}, []string{"3", "4"}, nil},
	}

	for _, testCase := range testCases {
		scope, resp := ctrl.validateScope(testCase.tr, testCase.grantData)
		assert.Equal(t, testCase.resp, resp)
		assert.Equal(t, testCase.scope, scope)
	}
}

func TestTokenCtrl_ValidateToken_WithPasswordGrant(t *testing.T) {
	ctrl, m := setupTokenCtrl()

	tr := &TokenRequest{
		Username:     "abc",
		Password:     "xyz",
		ClientID:     "1234",
		ClientSecret: "5678",
		GrantType:    PASSWORD,
		Scope:        "1 2",
	}

	td := &TokenData{}

	client := &DefaultClient{
		ClientID:  "1234",
		GrantType: []string{PASSWORD},
		Scope:     []string{"1", "2"},
	}

	m.store.On("GetUser", tr.Username, tr.Password).Return(tr.Username, nil, nil)
	m.store.On("GetClient", tr.ClientID).Return(client, nil)

	resp := ctrl.validateToken(tr, td)
	assert.Nil(t, resp)
	m.store.AssertExpectations(t)
}

func TestTokenCtrl_ValidateToken_WithClientGrant(t *testing.T) {
	ctrl, m := setupTokenCtrl()

	tr := &TokenRequest{
		ClientID:     "1234",
		ClientSecret: "5678",
		GrantType:    CLIENT_CREDENTIALS,
		Scope:        "1 2",
	}

	td := &TokenData{}

	client := &DefaultClient{
		ClientID:     "1234",
		ClientSecret: "5678",
		GrantType:    []string{CLIENT_CREDENTIALS},
		Scope:        []string{"1", "2"},
	}

	m.store.On("GetClient", tr.ClientID).Return(client, nil)

	resp := ctrl.validateToken(tr, td)
	assert.Nil(t, resp)
	m.store.AssertExpectations(t)
}

func TestTokenCtrl_HandleAccessToken(t *testing.T) {
	ctrl, m := setupTokenCtrl()

	tr := &TokenRequest{
		ClientID:     "1234",
		ClientSecret: "5678",
		GrantType:    CLIENT_CREDENTIALS,
		Scope:        "1 2",
	}

	client := &DefaultClient{
		ClientID:     "1234",
		ClientSecret: "5678",
		GrantType:    []string{CLIENT_CREDENTIALS},
		Scope:        []string{"1", "2"},
	}

	expResp := newRespData(nil)

	m.store.On("GetClient", tr.ClientID).Return(client, nil)
	m.tokenRespType.On("GetAccessToken", mock.Anything, false).Return(expResp)

	resp := ctrl.handleAccessToken(tr)
	assert.Equal(t, 200, resp.code)
	assert.False(t, resp.IsError())
}
