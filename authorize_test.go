package clover

import (
	"errors"
	"github.com/plimble/unik"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockAuthCtrl struct {
	store *Mockallstore
}

func setupAuthCtrl() (*authController, *mockAuthCtrl) {
	store := NewMockallstore()
	config := NewAuthConfig(store)
	config.AddAuthCodeGrant(store)
	authRespType := newCodeRespType(config, unik.NewUUIDV1())
	tokenRespType := newTokenRespType(config, unik.NewUUID1Base64())
	mock := &mockAuthCtrl{store}
	ctrl := newAuthController(config, authRespType, tokenRespType)

	return ctrl, mock
}

func TestAuthCtrl_ValidateClientID(t *testing.T) {
	ctrl, m := setupAuthCtrl()

	ar := &authorizeRequest{clientID: "123"}
	m.store.On("GetClient", ar.clientID).Return(&DefaultClient{}, nil)
	resp := ctrl.validateClientID(ar)
	assert.Nil(t, resp)
	assert.NotNil(t, ar.client)
	m.store.AssertExpectations(t)
}

func TestAuthCtrl_ValidateClientID_WithEmptyClientID(t *testing.T) {
	ctrl, m := setupAuthCtrl()

	ar := &authorizeRequest{}
	resp := ctrl.validateClientID(ar)
	assert.Equal(t, errNoClientID, resp)
	assert.Nil(t, ar.client)
	m.store.AssertExpectations(t)
}

func TestAuthCtrl_ValidateClientID_WithInvalidClient(t *testing.T) {
	ctrl, m := setupAuthCtrl()

	ar := &authorizeRequest{clientID: "123"}
	m.store.On("GetClient", ar.clientID).Return(nil, errors.New("not found"))
	resp := ctrl.validateClientID(ar)
	assert.Equal(t, errInvalidClientID, resp)
	assert.Nil(t, ar.client)
	m.store.AssertExpectations(t)
}

func TestAuthCtrl_ValidateState(t *testing.T) {
	ctrl, _ := setupAuthCtrl()

	ctrl.config.StateParamRequired = true
	ar := &authorizeRequest{state: "1234"}

	//WithRequired
	resp := ctrl.validateState(ar)
	assert.Nil(t, resp)
	assert.Equal(t, "1234", ar.state)

	//WithRequiredAndEmpty
	ar.state = ""

	resp = ctrl.validateState(ar)
	assert.Equal(t, errStateRequired, resp)
	assert.Equal(t, "", ar.state)
}

func TestAuthCtrl_ValidateRedirectURI(t *testing.T) {
	ctrl, _ := setupAuthCtrl()

	testCases := []struct {
		ar          *authorizeRequest
		resp        *Response
		redirectURI string
	}{
		{&authorizeRequest{redirectURI: "123", client: &DefaultClient{RedirectURI: "123"}}, nil, "123"},
		{&authorizeRequest{redirectURI: "", client: &DefaultClient{RedirectURI: "123"}}, nil, "123"},
		{&authorizeRequest{redirectURI: "321", client: &DefaultClient{RedirectURI: "123"}}, errRedirectMismatch, "321"},
		{&authorizeRequest{redirectURI: "", client: &DefaultClient{RedirectURI: ""}}, errNoRedirectURI, ""},
	}

	for _, testCase := range testCases {
		resp := ctrl.validateRedirectURI(testCase.ar)
		assert.Equal(t, testCase.resp, resp)
		assert.Equal(t, testCase.redirectURI, testCase.ar.redirectURI)
	}
}

func TestAuthCtrl_ValidateRespType_WithImplicit(t *testing.T) {
	ctrl, _ := setupAuthCtrl()
	ctrl.config.AllowImplicit = true

	testCases := []struct {
		ar       *authorizeRequest
		resp     *Response
		respType AuthRespType
	}{
		{&authorizeRequest{responseType: "", client: &DefaultClient{GrantType: []string{}}}, errInvalidRespType, nil},
		{&authorizeRequest{responseType: RESP_TYPE_CODE, client: &DefaultClient{GrantType: []string{AUTHORIZATION_CODE}}}, nil, ctrl.authRespType},
		{&authorizeRequest{responseType: RESP_TYPE_CODE, client: &DefaultClient{GrantType: []string{}}}, errUnAuthorizedGrant, nil},
		{&authorizeRequest{responseType: RESP_TYPE_TOKEN, client: &DefaultClient{GrantType: []string{IMPLICIT}}}, nil, ctrl.tokenRespType},
		{&authorizeRequest{responseType: RESP_TYPE_TOKEN, client: &DefaultClient{GrantType: []string{}}}, errUnAuthorizedGrant, nil},
		{&authorizeRequest{responseType: "xxxxx", client: &DefaultClient{GrantType: []string{}}}, errCodeUnSupportedGrant, nil},
	}

	for _, testCase := range testCases {
		resp := ctrl.validateRespType(testCase.ar)
		assert.Equal(t, testCase.resp, resp)
		assert.Equal(t, testCase.respType, testCase.ar.respType)
	}
}

func TestAuthCtrl_ValidateRespType_WithDisableImplicit(t *testing.T) {
	ctrl, _ := setupAuthCtrl()
	ctrl.config.AllowImplicit = false

	ar := &authorizeRequest{
		responseType: RESP_TYPE_TOKEN,
		client: &DefaultClient{
			GrantType: []string{IMPLICIT},
		},
	}

	resp := ctrl.validateRespType(ar)
	assert.Equal(t, errUnSupportedImplicit, resp)
	assert.Nil(t, ar.respType)
}

func TestAuthCtrl_ValidateScope(t *testing.T) {
	ctrl, _ := setupAuthCtrl()
	ctrl.config.DefaultScopes = []string{"a", "b", "c"}

	testCases := []struct {
		ar    *authorizeRequest
		resp  *Response
		scope []string
	}{
		{&authorizeRequest{scope: "", client: &DefaultClient{Scope: []string{}}}, nil, ctrl.config.DefaultScopes},
		{&authorizeRequest{scope: "1", client: &DefaultClient{Scope: []string{}}}, errUnSupportedScope, nil},
		{&authorizeRequest{scope: "a 1", client: &DefaultClient{Scope: []string{"a 2"}}}, errUnSupportedScope, nil},
		{&authorizeRequest{scope: "1 2", client: &DefaultClient{Scope: []string{"1", "2"}}}, nil, []string{"1", "2"}},
	}

	for _, testCase := range testCases {
		resp := ctrl.validateScope(testCase.ar)
		assert.Equal(t, testCase.resp, resp)
		assert.Equal(t, testCase.scope, testCase.ar.scopeArr)
	}
}

func TestAuthCtrl_ValidateAuthRequest(t *testing.T) {
	ctrl, m := setupAuthCtrl()

	ar := &authorizeRequest{
		clientID:     "123",
		responseType: "code",
	}

	client := &DefaultClient{
		RedirectURI: "http://localhost",
		GrantType:   []string{AUTHORIZATION_CODE},
	}

	m.store.On("GetClient", ar.clientID).Return(client, nil)

	resp := ctrl.validateAuthRequest(ar)
	m.store.AssertExpectations(t)
	assert.Nil(t, resp)
}

func TestAuthCtrl_ValidateAuthRequest_WithErrGrantImplicitRedirect(t *testing.T) {
	ctrl, m := setupAuthCtrl()
	ctrl.config.AllowImplicit = true

	ar := &authorizeRequest{
		clientID:     "123",
		responseType: "token",
		state:        "333",
	}

	client := &DefaultClient{
		RedirectURI: "http://localhost",
		GrantType:   []string{},
	}

	m.store.On("GetClient", ar.clientID).Return(client, nil)

	resp := ctrl.validateAuthRequest(ar)
	expResp := errUnAuthorizedGrant.clone().setRedirect("http://localhost", RESP_TYPE_TOKEN, "333")
	m.store.AssertExpectations(t)
	assert.Equal(t, expResp, resp)
}

func TestAuthCtrl_ValidateAuthRequest_WithErrScopeCodeRedirect(t *testing.T) {
	ctrl, m := setupAuthCtrl()

	ar := &authorizeRequest{
		clientID:     "123",
		responseType: "code",
		state:        "333",
		scope:        "1",
	}

	client := &DefaultClient{
		RedirectURI: "http://localhost",
		GrantType:   []string{AUTHORIZATION_CODE},
		Scope:       []string{"a"},
	}

	m.store.On("GetClient", ar.clientID).Return(client, nil)

	resp := ctrl.validateAuthRequest(ar)
	expResp := errUnSupportedScope.clone().setRedirect("http://localhost", RESP_TYPE_CODE, "333")
	m.store.AssertExpectations(t)
	assert.Equal(t, expResp, resp)
}

func TestHandleAuthorize_WithCode(t *testing.T) {
	ctrl, m := setupAuthCtrl()

	ar := &authorizeRequest{
		clientID:     "123",
		responseType: "code",
		state:        "333",
		scope:        "1",
	}

	client := &DefaultClient{
		RedirectURI: "http://localhost",
		GrantType:   []string{AUTHORIZATION_CODE},
		Scope:       []string{"1"},
	}

	m.store.On("GetClient", ar.clientID).Return(client, nil)
	m.store.On("SetAuthorizeCode", mock.Anything).Return(nil)

	resp := ctrl.handleAuthorize(ar, true)
	m.store.AssertExpectations(t)
	assert.False(t, resp.IsError())
	assert.True(t, resp.IsRedirect())
	assert.False(t, resp.isFragment)
}

func TestHandleAuthorize_WithImplicit(t *testing.T) {
	ctrl, m := setupAuthCtrl()
	ctrl.config.AllowImplicit = true

	ar := &authorizeRequest{
		clientID:     "123",
		responseType: "token",
		state:        "333",
		scope:        "1",
	}

	client := &DefaultClient{
		RedirectURI: "http://localhost",
		GrantType:   []string{IMPLICIT},
		Scope:       []string{"1"},
	}

	m.store.On("GetClient", ar.clientID).Return(client, nil)
	m.store.On("SetAccessToken", mock.Anything).Return(nil)

	resp := ctrl.handleAuthorize(ar, true)
	m.store.AssertExpectations(t)
	assert.False(t, resp.IsError())
	assert.True(t, resp.IsRedirect())
	assert.True(t, resp.isFragment)
}

func TestHandleAuthorize_WithUserDeclined(t *testing.T) {
	ctrl, _ := setupAuthCtrl()

	ar := &authorizeRequest{
		responseType: "token",
	}

	resp := ctrl.handleAuthorize(ar, false)
	expResp := errUserDeniedAccess.clone().setRedirect("", RESP_TYPE_TOKEN, "")
	assert.Equal(t, expResp, resp)
}
