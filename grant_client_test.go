package clover

import (
	"github.com/plimble/utils/errors2"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GrantClientSuite struct {
	suite.Suite
	store *Mockallstore
	grant GrantType
}

func TestGrantClientSuite(t *testing.T) {
	suite.Run(t, &GrantClientSuite{})
}

func (t *GrantClientSuite) SetupTest() {
	t.store = NewMockallstore()
	t.grant = NewClientCredential(t.store)
}

func (t *GrantClientSuite) TestValidate() {
	client := &TestClient{
		ClientID:     "123",
		UserID:       "abc",
		ClientSecret: "321",
		Scope:        []string{"1", "2"},
		GrantType:    []string{CLIENT_CREDENTIALS},
	}

	tr := &TokenRequest{
		ClientID:     "123",
		ClientSecret: "321",
	}

	expGrantData := &GrantData{
		ClientID: client.ClientID,
		UserID:   client.UserID,
		Scope:    client.Scope,
	}

	t.store.On("GetClient", tr.ClientID).Return(client, nil)

	grantData, resp := t.grant.Validate(tr)
	t.store.AssertExpectations(t.T())
	t.Nil(resp)
	t.EqualValues(grantData, expGrantData)
}

func (t *GrantClientSuite) TestValidate_WithNotFoundClient() {
	tr := &TokenRequest{
		ClientID:     "123",
		ClientSecret: "321",
	}

	t.store.On("GetClient", tr.ClientID).Return(nil, errors2.NewAnyError())

	grantData, resp := t.grant.Validate(tr)
	t.store.AssertExpectations(t.T())
	t.Equal(errInvalidClientCredentail, resp)
	t.Nil(grantData)
}

func (t *GrantClientSuite) TestValidate_WithIvalidClientSecret() {
	client := &TestClient{
		ClientID:     "123",
		UserID:       "abc",
		ClientSecret: "321",
		Scope:        []string{"1", "2"},
		GrantType:    []string{CLIENT_CREDENTIALS},
	}

	tr := &TokenRequest{
		ClientID:     "123",
		ClientSecret: "xxxx",
	}

	t.store.On("GetClient", tr.ClientID).Return(client, nil)

	grantData, resp := t.grant.Validate(tr)
	t.store.AssertExpectations(t.T())
	t.Equal(errInvalidClientCredentail, resp)
	t.Nil(grantData)
}

func (t *GrantClientSuite) TestName() {
	t.Equal(CLIENT_CREDENTIALS, t.grant.Name())
}

func (t *GrantClientSuite) TestIncludeRefreshToken() {
	t.False(t.grant.IncludeRefreshToken())
}

func (t *GrantClientSuite) TestBeforeCreateAccessToken() {
	tr := &TokenRequest{}
	td := &TokenData{}
	t.Nil(t.grant.BeforeCreateAccessToken(tr, td))
}
