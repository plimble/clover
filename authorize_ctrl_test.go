package clover

import (
	"github.com/plimble/utils/errors2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuthorizeCtrlSuite struct {
	suite.Suite
	store *Mockallstore
	ctrl  *authorizeCtrl
}

func TestAuthorizeCtrlSuite(t *testing.T) {
	suite.Run(t, &AuthorizeCtrlSuite{})
}

func (t *AuthorizeCtrlSuite) SetupTest() {
	config := &AuthServerConfig{}
	t.store = NewMockallstore()
	t.ctrl = newAuthorizeCtrl(t.store, config)
}

func (t *AuthorizeCtrlSuite) TestValidateClientID() {
	expc := genTestClient()
	ar := &authorizeRequest{clientID: expc.ClientID}

	t.store.On("GetClient", ar.clientID).Return(expc, nil)

	client, resp := t.ctrl.validateClientID(ar)

	t.Nil(resp)
	t.Equal(expc, client)
	t.store.AssertExpectations(t.T())
}

func (t *AuthorizeCtrlSuite) TestValidateClientID_WithEmptyClientID() {
	ar := &authorizeRequest{}
	client, resp := t.ctrl.validateClientID(ar)

	t.Equal(errNoClientID, resp)
	t.Nil(client)
}

func (t *AuthorizeCtrlSuite) TestValidateClientID_WithInvalidClient() {
	ar := &authorizeRequest{clientID: "123"}
	t.store.On("GetClient", ar.clientID).Return(nil, errors2.NewAnyError())

	client, resp := t.ctrl.validateClientID(ar)

	t.Equal(errInvalidClientID, resp)
	t.Nil(client)
	t.store.AssertExpectations(t.T())
}

func (t *AuthorizeCtrlSuite) TestValidateState() {
	t.ctrl.config.StateParamRequired = true
	ar := &authorizeRequest{state: "1234"}

	//WithRequired
	state, resp := t.ctrl.validateState(ar)
	t.Nil(resp)
	t.Equal("1234", state)

	//WithRequiredAndEmpty
	ar.state = ""

	state, resp = t.ctrl.validateState(ar)
	t.Equal(errStateRequired, resp)
	t.Equal("", state)
}

func (t *AuthorizeCtrlSuite) TestValidateRedirectURI() {
	testCases := []struct {
		ar          *authorizeRequest
		client      Client
		resp        *Response
		redirectURI string
	}{
		{&authorizeRequest{redirectURI: "123"}, &TestClient{RedirectURI: "123"}, nil, "123"},
		{&authorizeRequest{redirectURI: ""}, &TestClient{RedirectURI: "123"}, nil, "123"},
		{&authorizeRequest{redirectURI: "321"}, &TestClient{RedirectURI: "123"}, errRedirectMismatch, ""},
		{&authorizeRequest{redirectURI: ""}, &TestClient{RedirectURI: ""}, errNoRedirectURI, ""},
	}

	for _, testCase := range testCases {
		redirectURI, resp := t.ctrl.validateRedirectURI(testCase.client, testCase.ar)
		t.Equal(testCase.resp, resp)
		t.Equal(testCase.redirectURI, redirectURI)
	}
}

func (t *AuthorizeCtrlSuite) TestValidateRespType() {
	mockResp := NewMockAuthorizeRespType()
	authRespTypes := map[string]AuthorizeRespType{IMPLICIT: mockResp}
	ar := &authorizeRequest{responseType: IMPLICIT}
	client := genTestClient()
	client.GrantType = []string{IMPLICIT}

	mockResp.On("SupportGrant").Return(IMPLICIT)

	respType, resp := t.ctrl.validateRespType(client, ar, authRespTypes)
	t.Equal(mockResp, respType)
	t.Nil(resp)
	mockResp.AssertExpectations(t.T())
}

func (t *AuthorizeCtrlSuite) TestValidateRespType_WithNoRespType() {
	authRespTypes := map[string]AuthorizeRespType{}
	ar := &authorizeRequest{responseType: "code"}
	client := genTestClient()

	respType, resp := t.ctrl.validateRespType(client, ar, authRespTypes)
	t.Nil(respType)
	t.Equal(errInvalidRespType, resp)
}

func (t *AuthorizeCtrlSuite) TestValidateRespType_WithUnsupportGrant() {
	mockResp := NewMockAuthorizeRespType()
	authRespTypes := map[string]AuthorizeRespType{"token": mockResp}
	ar := &authorizeRequest{responseType: "token"}

	client := genTestClient()
	client.GrantType = []string{AUTHORIZATION_CODE}

	mockResp.On("SupportGrant").Return(IMPLICIT)
	respType, resp := t.ctrl.validateRespType(client, ar, authRespTypes)
	t.Nil(respType)
	t.Equal(errUnAuthorizedGrant, resp)
	mockResp.AssertExpectations(t.T())
}

func (t *AuthorizeCtrlSuite) TestValidateScope() {
	t.ctrl.config.DefaultScopes = []string{"a", "b", "c"}

	testCases := []struct {
		ar     *authorizeRequest
		client Client
		resp   *Response
		scope  []string
	}{
		{&authorizeRequest{scope: ""}, &TestClient{Scope: []string{}}, nil, t.ctrl.config.DefaultScopes},
		{&authorizeRequest{scope: "1"}, &TestClient{Scope: []string{}}, errUnSupportedScope, nil},
		{&authorizeRequest{scope: "a 1"}, &TestClient{Scope: []string{"a 2"}}, errUnSupportedScope, nil},
		{&authorizeRequest{scope: "1 2"}, &TestClient{Scope: []string{"1", "2"}}, nil, []string{"1", "2"}},
	}

	for _, testCase := range testCases {
		scope, resp := t.ctrl.validateScope(testCase.client, testCase.ar)
		t.Equal(testCase.resp, resp)
		t.Equal(testCase.scope, scope)
	}
}

func (t *AuthorizeCtrlSuite) TestValidate() {
	mockResp := NewMockAuthorizeRespType()
	authRespTypes := map[string]AuthorizeRespType{"code": mockResp}
	ar := &authorizeRequest{
		clientID:     "123",
		responseType: "code",
	}

	client := genTestClient()
	client.RedirectURI = "http://localhost"
	client.GrantType = []string{AUTHORIZATION_CODE}

	mockResp.On("SupportGrant").Return(AUTHORIZATION_CODE)
	t.store.On("GetClient", ar.clientID).Return(client, nil)

	ad, resp := t.ctrl.validate(ar, authRespTypes)
	t.store.AssertExpectations(t.T())
	mockResp.AssertExpectations(t.T())
	t.Nil(resp)
	t.NotNil(ad)
}

func (t *AuthorizeCtrlSuite) TestValidateAuthRequest_WithErrScopeCodeRedirect() {
	mockResp := NewMockAuthorizeRespType()
	authRespTypes := map[string]AuthorizeRespType{"code": mockResp}

	ar := &authorizeRequest{
		clientID:     "123",
		responseType: "code",
		state:        "333",
		scope:        "1",
	}

	client := genTestClient()
	client.RedirectURI = "http://localhost"
	client.GrantType = []string{AUTHORIZATION_CODE}
	client.Scope = []string{"a"}

	mockResp.On("SupportGrant").Return(AUTHORIZATION_CODE)
	mockResp.On("IsImplicit").Return(false)
	t.store.On("GetClient", ar.clientID).Return(client, nil)

	ad, resp := t.ctrl.validate(ar, authRespTypes)
	expResp := errUnSupportedScope.setRedirect("http://localhost", false, "333")
	t.store.AssertExpectations(t.T())
	mockResp.AssertExpectations(t.T())
	t.Equal(expResp, resp)
	t.Nil(ad)
}

func (t *AuthorizeCtrlSuite) TestValidateAuthRequest_WithErrScopeImplicitRedirect() {
	mockResp := NewMockAuthorizeRespType()
	authRespTypes := map[string]AuthorizeRespType{"code": mockResp}

	ar := &authorizeRequest{
		clientID:     "123",
		responseType: "code",
		state:        "333",
		scope:        "1",
	}

	client := genTestClient()
	client.RedirectURI = "http://localhost"
	client.GrantType = []string{AUTHORIZATION_CODE}
	client.Scope = []string{"a"}

	mockResp.On("SupportGrant").Return(AUTHORIZATION_CODE)
	mockResp.On("IsImplicit").Return(true)
	t.store.On("GetClient", ar.clientID).Return(client, nil)

	ad, resp := t.ctrl.validate(ar, authRespTypes)
	expResp := errUnSupportedScope.setRedirect("http://localhost", true, "333")
	t.store.AssertExpectations(t.T())
	mockResp.AssertExpectations(t.T())
	t.Equal(expResp, resp)
	t.Nil(ad)
}

func (t *AuthorizeCtrlSuite) TestHandleAuthorize() {
	mockResp := NewMockAuthorizeRespType()
	authRespTypes := map[string]AuthorizeRespType{"code": mockResp}

	ar := &authorizeRequest{
		clientID:     "123",
		responseType: "code",
		state:        "333",
		scope:        "1",
	}

	client := genTestClient()
	client.RedirectURI = "http://localhost"
	client.GrantType = []string{AUTHORIZATION_CODE}
	client.Scope = []string{"1"}

	mockResp.On("SupportGrant").Return(AUTHORIZATION_CODE)
	mockResp.On("Response", mock.Anything).Return(newRespData(nil).setRedirect("", false, ""))
	t.store.On("GetClient", ar.clientID).Return(client, nil)

	resp := t.ctrl.authorize(ar, authRespTypes, true)
	t.store.AssertExpectations(t.T())
	mockResp.AssertExpectations(t.T())
	t.False(resp.IsError())
	t.True(resp.IsRedirect())
	t.False(resp.isFragment)
}

func (t *AuthorizeCtrlSuite) TestHandleAuthorize_WithUserDeclined() {
	authRespTypes := map[string]AuthorizeRespType{}
	ar := &authorizeRequest{
		responseType: "token",
		redirectURI:  "http://localhost",
		state:        "123",
	}

	resp := t.ctrl.authorize(ar, authRespTypes, false)
	expResp := errUserDeniedAccess.setRedirect("http://localhost", false, "123")
	t.Equal(expResp, resp)
}
