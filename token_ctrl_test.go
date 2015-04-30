package clover

import (
	"github.com/plimble/utils/errors2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TokenCtrlSuite struct {
	suite.Suite
	ctrl  *tokenCtrl
	store *Mockallstore
}

func TestTokenCtrlSuite(t *testing.T) {
	suite.Run(t, &TokenCtrlSuite{})
}

func (t *TokenCtrlSuite) SetupTest() {
	config := &AuthServerConfig{}
	t.store = NewMockallstore()
	t.ctrl = newTokenCtrl(t.store, t.store, config)
}

func (t *TokenCtrlSuite) TestValidateGrantType() {
	mockGrant := NewMockGrantType()
	grants := map[string]GrantType{CLIENT_CREDENTIALS: mockGrant}

	tr := &TokenRequest{GrantType: CLIENT_CREDENTIALS}

	grant, resp := t.ctrl.validateGrantType(tr, grants)
	mockGrant.AssertExpectations(t.T())
	t.Nil(resp)
	t.Equal(mockGrant, grant)
}

func (t *TokenCtrlSuite) TestValidateGrantType_WithGrantRequired() {
	mockGrant := NewMockGrantType()
	grants := map[string]GrantType{CLIENT_CREDENTIALS: mockGrant}

	tr := &TokenRequest{}

	grant, resp := t.ctrl.validateGrantType(tr, grants)
	mockGrant.AssertExpectations(t.T())
	t.Equal(errGrantTypeRequired, resp)
	t.Nil(grant)
}

func (t *TokenCtrlSuite) TestValidateGrantType_WithUnsupportedGrant() {
	mockGrant := NewMockGrantType()
	grants := map[string]GrantType{CLIENT_CREDENTIALS: mockGrant}

	tr := &TokenRequest{GrantType: AUTHORIZATION_CODE}

	grant, resp := t.ctrl.validateGrantType(tr, grants)
	mockGrant.AssertExpectations(t.T())
	t.Equal(errUnSupportedGrantType, resp)
	t.Nil(grant)
}

func (t *TokenCtrlSuite) TestValidateClient() {
	expc := genTestClient()

	tr := &TokenRequest{
		ClientID: expc.GetClientID(),
	}

	grantData := &GrantData{
		ClientID: expc.GetClientID(),
	}

	t.store.On("GetClient", tr.ClientID).Return(expc, nil)

	client, resp := t.ctrl.validateClient(tr, grantData)
	t.store.AssertExpectations(t.T())
	t.Nil(resp)
	t.Equal(expc, client)
}

func (t *TokenCtrlSuite) TestValidateClient_WithInvalidClient() {
	expc := genTestClient()

	tr := &TokenRequest{
		ClientID: expc.GetClientID(),
	}

	grantData := &GrantData{
		ClientID: expc.GetClientID(),
	}

	t.store.On("GetClient", tr.ClientID).Return(nil, errors2.NewAnyError())

	client, resp := t.ctrl.validateClient(tr, grantData)
	t.store.AssertExpectations(t.T())
	t.Nil(client)
	t.Equal(errInvalidClientID, resp)
}

func (t *TokenCtrlSuite) TestValidateClient_WithInvalidCredential() {
	expc := genTestClient()

	tr := &TokenRequest{
		ClientID: expc.GetClientID(),
	}

	grantData := &GrantData{
		ClientID: "11",
	}

	t.store.On("GetClient", tr.ClientID).Return(expc, nil)

	client, resp := t.ctrl.validateClient(tr, grantData)
	t.store.AssertExpectations(t.T())
	t.Nil(client)
	t.Equal(errInvalidClientCredentail, resp)
}

func (t *TokenCtrlSuite) TestValidateScope_CheckGrantScope() {
	tr := &TokenRequest{Scope: "1 2"}
	grantData := &GrantData{Scope: []string{"1", "2", "3"}}
	client := &TestClient{}

	scope, resp := t.ctrl.validateScope(tr, grantData, client)
	t.Nil(resp)
	t.Equal([]string{"1", "2"}, scope)
}

func (t *TokenCtrlSuite) TestValidateScope_CheckGrantScopeFailed() {
	tr := &TokenRequest{Scope: "1 2"}
	grantData := &GrantData{Scope: []string{"1"}}
	client := &TestClient{}

	scope, resp := t.ctrl.validateScope(tr, grantData, client)
	t.Equal(errInvalidScopeRequest, resp)
	t.Nil(scope)
}

func (t *TokenCtrlSuite) TestValidateScope_CheckClientScope() {
	tr := &TokenRequest{Scope: "1 2"}
	grantData := &GrantData{}
	client := &TestClient{Scope: []string{"1", "2", "3"}}

	scope, resp := t.ctrl.validateScope(tr, grantData, client)
	t.Nil(resp)
	t.Equal([]string{"1", "2"}, scope)
}

func (t *TokenCtrlSuite) TestValidateScope_CheckClientScopeFailed() {
	tr := &TokenRequest{Scope: "1 2 4"}
	grantData := &GrantData{}
	client := &TestClient{Scope: []string{"1", "2", "3"}}

	scope, resp := t.ctrl.validateScope(tr, grantData, client)
	t.Equal(errInvalidScopeRequest, resp)
	t.Nil(scope)
}

func (t *TokenCtrlSuite) TestValidateScope_ExistScope() {
	tr := &TokenRequest{Scope: "1 2"}
	grantData := &GrantData{}
	client := &TestClient{}

	t.store.On("ExistScopes", []string{"1", "2"}).Return(true, nil)

	scope, resp := t.ctrl.validateScope(tr, grantData, client)
	t.Nil(resp)
	t.Equal([]string{"1", "2"}, scope)
}

func (t *TokenCtrlSuite) TestValidateScope_ExistScopeFailed() {
	tr := &TokenRequest{Scope: "1 2"}
	grantData := &GrantData{}
	client := &TestClient{}

	t.store.On("ExistScopes", []string{"1", "2"}).Return(false, nil)

	scope, resp := t.ctrl.validateScope(tr, grantData, client)
	t.Nil(scope)
	t.Equal(errUnSupportedScope, resp)
}

func (t *TokenCtrlSuite) TestValidateScope_NoReqScopeUseGrantScope() {
	tr := &TokenRequest{}
	grantData := &GrantData{Scope: []string{"1", "2"}}
	client := &TestClient{}

	scope, resp := t.ctrl.validateScope(tr, grantData, client)
	t.Equal([]string{"1", "2"}, scope)
	t.Nil(resp)
}

func (t *TokenCtrlSuite) TestValidateScope_DefaultScope() {
	tr := &TokenRequest{}
	grantData := &GrantData{}
	client := &TestClient{}

	t.store.On("GetDefaultScope", client.GetClientID()).Return([]string{"1", "2"}, nil)

	scope, resp := t.ctrl.validateScope(tr, grantData, client)
	t.Equal([]string{"1", "2"}, scope)
	t.Nil(resp)
}

func (t *TokenCtrlSuite) TestValidateScope_EmptyDefaultScope() {
	tr := &TokenRequest{}
	grantData := &GrantData{}
	client := &TestClient{}

	t.store.On("GetDefaultScope", client.GetClientID()).Return(nil, nil)

	scope, resp := t.ctrl.validateScope(tr, grantData, client)
	t.Nil(scope)
	t.Equal(errNoScope, resp)
}

func (t *TokenCtrlSuite) TestToken() {
	mockRespType := NewMockAccessTokenRespType()
	mockGrant := NewMockGrantType()
	grants := map[string]GrantType{PASSWORD: mockGrant}

	tr := &TokenRequest{
		Username:     "abc",
		Password:     "xyz",
		ClientID:     "1234",
		ClientSecret: "5678",
		GrantType:    PASSWORD,
		Scope:        "1 2",
	}

	grantData := &GrantData{
		UserID: "user1",
	}

	client := &TestClient{
		ClientID:  "1234",
		GrantType: []string{PASSWORD},
		Scope:     []string{"1", "2"},
	}

	mockGrant.On("Validate", tr).Return(grantData, nil)
	mockGrant.On("Name").Return(PASSWORD)
	mockGrant.On("BeforeCreateAccessToken", tr, mock.Anything).Return(nil)
	mockGrant.On("IncludeRefreshToken").Return(false)

	mockRespType.On("Response", mock.Anything, mock.Anything).Return(newRespData(nil))

	t.store.On("GetClient", tr.ClientID).Return(client, nil)

	resp := t.ctrl.token(tr, mockRespType, grants)
	t.Equal(newRespData(nil), resp)
	t.store.AssertExpectations(t.T())
	mockRespType.AssertExpectations(t.T())
	mockGrant.AssertExpectations(t.T())
}

func (t *TokenCtrlSuite) TestToken_WithInvalidGrant() {
	mockRespType := NewMockAccessTokenRespType()
	mockGrant := NewMockGrantType()
	grants := map[string]GrantType{PASSWORD: mockGrant}

	tr := &TokenRequest{
		Username:     "abc",
		Password:     "xyz",
		ClientID:     "1234",
		ClientSecret: "5678",
		GrantType:    PASSWORD,
		Scope:        "1 2",
	}

	grantData := &GrantData{
		UserID: "user1",
	}

	client := &TestClient{
		ClientID:  "1234",
		GrantType: []string{CLIENT_CREDENTIALS},
		Scope:     []string{"1", "2"},
	}

	mockGrant.On("Validate", tr).Return(grantData, nil)
	mockGrant.On("Name").Return(PASSWORD)

	t.store.On("GetClient", tr.ClientID).Return(client, nil)

	resp := t.ctrl.token(tr, mockRespType, grants)
	t.Equal(errUnAuthorizedGrant, resp)
	t.store.AssertExpectations(t.T())
	mockRespType.AssertExpectations(t.T())
	mockGrant.AssertExpectations(t.T())
}
