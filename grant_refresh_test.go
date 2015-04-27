package clover

import (
	"github.com/plimble/utils/errors2"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GrantRefreshSuite struct {
	suite.Suite
	store *Mockallstore
	grant GrantType
}

func TestGrantRefreshSuite(t *testing.T) {
	suite.Run(t, &GrantRefreshSuite{})
}

func (t *GrantRefreshSuite) SetupTest() {
	t.store = NewMockallstore()
	t.grant = NewRefreshToken(t.store)
}

func (t *GrantRefreshSuite) TestValidate() {
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
		ClientID: rt.ClientID,
		UserID:   rt.UserID,
		Scope:    rt.Scope,
	}

	t.store.On("GetRefreshToken", tr.RefreshToken).Return(rt, nil)

	grantData, resp := t.grant.Validate(tr)
	t.Nil(resp)
	t.EqualValues(expGrantData, grantData)
}

func (t *GrantRefreshSuite) TestValidate_WithNotFound() {
	tr := &TokenRequest{
		RefreshToken: "xyz",
	}

	t.store.On("GetRefreshToken", tr.RefreshToken).Return(nil, errors2.NewAnyError())

	grantData, resp := t.grant.Validate(tr)
	t.Equal(errInvalidRefreshToken, resp)
	t.Nil(grantData)
}

func (t *GrantRefreshSuite) TestValidate_WithTokenExpired() {
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

	t.store.On("GetRefreshToken", tr.RefreshToken).Return(rt, nil)

	grantData, resp := t.grant.Validate(tr)
	t.Equal(errRefreshTokenExpired, resp)
	t.Nil(grantData)
}

func (t *GrantRefreshSuite) TestValidate_WithTokenEmpty() {
	tr := &TokenRequest{
		RefreshToken: "",
	}

	grantData, resp := t.grant.Validate(tr)
	t.Equal(errRefreshTokenRequired, resp)
	t.Nil(grantData)
}

func (t *GrantRefreshSuite) TestName() {
	t.Equal(REFRESH_TOKEN, t.grant.Name())
}

func (t *GrantRefreshSuite) TestIncludeRefreshToken() {
	t.True(t.grant.IncludeRefreshToken())
}

func (t *GrantRefreshSuite) TestBeforeCreateAccessToken() {
	tr := &TokenRequest{
		RefreshToken: "1234",
	}
	td := &TokenData{}

	t.store.On("RemoveRefreshToken", tr.RefreshToken).Return(nil)
	t.Nil(t.grant.BeforeCreateAccessToken(tr, td))
}
