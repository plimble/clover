package clover

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockPasswordGrant struct {
	store *Mockallstore
}

func setUpPasswordGrant() (GrantType, *mockPasswordGrant) {
	store := NewMockallstore()
	grant := newPasswordGrant(store)
	mock := &mockPasswordGrant{store}
	return grant, mock
}

func TestPasswordGrant_Validate(t *testing.T) {
	c, m := setUpPasswordGrant()

	tr := &TokenRequest{
		Username: "abc",
		Password: "xyz",
	}

	expGrantData := &GrantData{
		UserID: "001",
		Scope:  []string{"1", "2", "3"},
	}

	m.store.On("GetUser", tr.Username, tr.Password).Return(expGrantData.UserID, expGrantData.Scope, nil)

	grantData, resp := c.Validate(tr)
	assert.Nil(t, resp)
	assert.EqualValues(t, grantData, expGrantData)
}

func TestPasswordGrant_GetGrantType(t *testing.T) {
	c, _ := setUpPasswordGrant()
	grant := c.GetGrantType()
	assert.Equal(t, PASSWORD, grant)
}

func TestPasswordGrant_CreateAccessToken(t *testing.T) {
	c, _ := setUpPasswordGrant()

	td := &TokenData{}
	respType := NewMockResponseType()
	respType.On("GetAccessToken", td, true).Return(nil)

	c.CreateAccessToken(td, respType)
}
