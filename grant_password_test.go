package clover

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
		Data:   map[string]interface{}{"a": 1, "b": "1"},
	}

	u := &DefaultUser{
		ID:       "001",
		Username: "test",
		Password: "1234",
		Scope:    []string{"1", "2", "3"},
		Data:     map[string]interface{}{"a": 1, "b": "1"},
	}

	m.store.On("GetUser", tr.Username, tr.Password).Return(u, nil)

	grantData, resp := c.Validate(tr)
	assert.Nil(t, resp)
	assert.EqualValues(t, grantData, expGrantData)
}

func TestPasswordGrant_Validate_WithNotFound(t *testing.T) {
	c, m := setUpPasswordGrant()

	tr := &TokenRequest{
		Username: "abc",
		Password: "xyz",
	}

	m.store.On("GetUser", tr.Username, tr.Password).Return(nil, errors.New("not found"))

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errInvalidUsernamePAssword, resp)
	assert.Nil(t, grantData)
}

func TestPasswordGrant_Validate_WithUsernamePasswordRequired(t *testing.T) {
	c, _ := setUpPasswordGrant()

	tr := &TokenRequest{
		Username: "",
		Password: "",
	}

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errUsernamePasswordRequired, resp)
	assert.Nil(t, grantData)
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

	respType.AssertExpectations(t)
}
