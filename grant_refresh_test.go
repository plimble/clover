package clover

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRefreshGrant struct {
	store *Mockallstore
}

func setUpRefreshGrant() (GrantType, *mockRefreshGrant) {
	store := NewMockallstore()
	grant := newRefreshGrant(store)
	mock := &mockRefreshGrant{store}
	return grant, mock
}

func TestRefreshGrant_Validate(t *testing.T) {
	c, m := setUpRefreshGrant()

	rt := &RefreshToken{
		RefreshToken: "xyz",
		ClientID:     "123",
		UserID:       "abc",
		Scope:        []string{"1", "2"},
		Expires:      addSecondUnix(60),
	}

	tr := &TokenRequest{
		RefreshToken: "xyz",
	}

	expGrantData := &GrantData{
		ClientID:     rt.ClientID,
		UserID:       rt.UserID,
		Scope:        rt.Scope,
		RefreshToken: rt.RefreshToken,
	}

	m.store.On("GetRefreshToken", tr.RefreshToken).Return(rt, nil)

	grantData, resp := c.Validate(tr)
	assert.Nil(t, resp)
	assert.EqualValues(t, grantData, expGrantData)
}

func TestRefreshGrant_Validate_WithNotFound(t *testing.T) {
	c, m := setUpRefreshGrant()

	tr := &TokenRequest{
		RefreshToken: "xyz",
	}

	m.store.On("GetRefreshToken", tr.RefreshToken).Return(nil, errors.New("not found"))

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errInvalidRefreshToken, resp)
	assert.Nil(t, grantData)
}

func TestRefreshGrant_Validate_WithTokenExpired(t *testing.T) {
	c, m := setUpRefreshGrant()

	rt := &RefreshToken{
		RefreshToken: "xyz",
		ClientID:     "123",
		UserID:       "abc",
		Scope:        []string{"1", "2"},
		Expires:      addSecondUnix(0),
	}

	tr := &TokenRequest{
		RefreshToken: "xyz",
	}

	m.store.On("GetRefreshToken", tr.RefreshToken).Return(rt, nil)

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errRefreshTokenExpired, resp)
	assert.Nil(t, grantData)
}

func TestRefreshGrant_Validate_WithTokenEmpty(t *testing.T) {
	c, _ := setUpRefreshGrant()

	tr := &TokenRequest{
		RefreshToken: "",
	}

	grantData, resp := c.Validate(tr)
	assert.Equal(t, errRefreshTokenRequired, resp)
	assert.Nil(t, grantData)

}

func TestRefreshGrant_GetGrantType(t *testing.T) {
	c, _ := setUpRefreshGrant()
	grant := c.GetGrantType()
	assert.Equal(t, REFRESH_TOKEN, grant)
}

func TestRefreshGrant_CreateAccessToken(t *testing.T) {
	c, m := setUpRefreshGrant()

	td := &TokenData{
		GrantData: &GrantData{
			RefreshToken: "123",
		},
	}
	respType := NewMockResponseType()

	m.store.On("RemoveRefreshToken", td.GrantData.RefreshToken).Return(nil)
	respType.On("GetAccessToken", td, true).Return(nil)

	c.CreateAccessToken(td, respType)

	m.store.AssertExpectations(t)
	respType.AssertExpectations(t)
}

func TestRefreshGrant_CreateAccessToken_WithCannotRemoveToken(t *testing.T) {
	c, m := setUpRefreshGrant()

	td := &TokenData{
		GrantData: &GrantData{
			RefreshToken: "123",
		},
	}
	respType := NewMockResponseType()

	m.store.On("RemoveRefreshToken", td.GrantData.RefreshToken).Return(errors.New("error"))
	resp := c.CreateAccessToken(td, respType)

	assert.True(t, resp.IsError())
	m.store.AssertExpectations(t)
	respType.AssertExpectations(t)
}
