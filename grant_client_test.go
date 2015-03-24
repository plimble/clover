package clover

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockClientGrant struct {
	store *Mockallstore
}

func setUpClientGrant() (GrantType, *mockClientGrant) {
	store := NewMockallstore()
	grant := newClientGrant(store)
	mock := &mockClientGrant{store}
	return grant, mock
}

func TestClientGrant_Validate(t *testing.T) {
	c, m := setUpClientGrant()

	client := &DefaultClient{
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
		ClientID:  client.ClientID,
		UserID:    client.UserID,
		Scope:     client.Scope,
		GrantType: client.GrantType,
	}

	m.store.On("GetClient", tr.ClientID).Return(client, nil)

	grantData, resp := c.Validate(tr)
	assert.Nil(t, resp)
	assert.EqualValues(t, grantData, expGrantData)
}

func TestClientGrant_Validate_WithNotFoundClient(t *testing.T) {
	c, m := setUpClientGrant()

	tr := &TokenRequest{
		ClientID:     "123",
		ClientSecret: "321",
	}

	m.store.On("GetClient", tr.ClientID).Return(nil, errors.New("not found"))

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errInvalidClientCredentail, resp)
	assert.Nil(t, grantData)
}

func TestClientGrant_Validate_WithIvalidClientSecret(t *testing.T) {
	c, m := setUpClientGrant()

	client := &DefaultClient{
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

	m.store.On("GetClient", tr.ClientID).Return(client, nil)

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errInvalidClientCredentail, resp)
	assert.Nil(t, grantData)
}

func TestClientGrant_GetGrantType(t *testing.T) {
	c, _ := setUpClientGrant()
	grant := c.GetGrantType()
	assert.Equal(t, CLIENT_CREDENTIALS, grant)
}

func TestClientGrant_CreateAccessToken(t *testing.T) {
	c, _ := setUpClientGrant()

	td := &TokenData{}
	respType := NewMockResponseType()
	respType.On("GetAccessToken", td, false).Return(nil)

	c.CreateAccessToken(td, respType)

	respType.AssertExpectations(t)
}
