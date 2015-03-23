package clover

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCodeGrant struct {
	store *Mockallstore
}

func setUpCodeGrant() (GrantType, *mockCodeGrant) {
	store := NewMockallstore()
	grant := newAuthCodeGrant(store)
	mock := &mockCodeGrant{store}
	return grant, mock
}

func TestCodeGrant_Validate(t *testing.T) {
	c, m := setUpCodeGrant()

	ac := &AuthorizeCode{
		ClientID:    "123",
		UserID:      "abc",
		Scope:       []string{"1", "2"},
		RedirectURI: "http://localhost",
		Expires:     addSecondUnix(60),
	}

	tr := &TokenRequest{
		Code:        "123",
		RedirectURI: "http://localhost",
	}

	expGrantData := &GrantData{
		ClientID: ac.ClientID,
		UserID:   ac.UserID,
		Scope:    ac.Scope,
	}

	m.store.On("GetAuthorizeCode", tr.Code).Return(ac, nil)

	grantData, resp := c.Validate(tr)
	assert.Nil(t, resp)
	assert.EqualValues(t, grantData, expGrantData)
}

func TestCodeGrant_Validate_WithNotFound(t *testing.T) {
	c, m := setUpCodeGrant()

	tr := &TokenRequest{
		Code:        "123",
		RedirectURI: "http://localhost",
	}

	m.store.On("GetAuthorizeCode", tr.Code).Return(nil, errors.New("not found"))

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errAuthCodeNotExist, resp)
	assert.Nil(t, grantData)
}

func TestCodeGrant_Validate_WithRedirectMisMatch(t *testing.T) {
	c, m := setUpCodeGrant()

	ac := &AuthorizeCode{
		ClientID:    "123",
		UserID:      "abc",
		Scope:       []string{"1", "2"},
		RedirectURI: "http://localhost2",
		Expires:     addSecondUnix(60),
	}

	tr := &TokenRequest{
		Code:        "123",
		RedirectURI: "http://localhost",
	}

	m.store.On("GetAuthorizeCode", tr.Code).Return(ac, nil)

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errRedirectMismatch, resp)
	assert.Nil(t, grantData)
}

func TestCodeGrant_Validate_WithCodeExpired(t *testing.T) {
	c, m := setUpCodeGrant()

	ac := &AuthorizeCode{
		ClientID:    "123",
		UserID:      "abc",
		Scope:       []string{"1", "2"},
		RedirectURI: "http://localhost",
		Expires:     0,
	}

	tr := &TokenRequest{
		Code:        "123",
		RedirectURI: "http://localhost",
	}

	m.store.On("GetAuthorizeCode", tr.Code).Return(ac, nil)

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errAuthCodeExpired, resp)
	assert.Nil(t, grantData)
}

func TestCodeGrant_Validate_WithCodeEmpty(t *testing.T) {
	c, _ := setUpCodeGrant()

	tr := &TokenRequest{
		Code:        "",
		RedirectURI: "http://localhost",
	}

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errCodeRequired, resp)
	assert.Nil(t, grantData)

}

func TestCodeGrant_GetGrantType(t *testing.T) {
	c, _ := setUpCodeGrant()
	grant := c.GetGrantType()
	assert.Equal(t, AUTHORIZATION_CODE, grant)
}

func TestCodeGrant_CreateAccessToken(t *testing.T) {
	c, _ := setUpCodeGrant()

	td := &TokenData{}
	respType := NewMockResponseType()
	respType.On("GetAccessToken", td, true).Return(nil)

	c.CreateAccessToken(td, respType)
}
