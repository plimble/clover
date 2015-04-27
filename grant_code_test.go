package clover

import (
	"github.com/plimble/utils/errors2"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GrantCodeSuite struct {
	suite.Suite
	store *Mockallstore
	grant GrantType
}

func TestGrantCodeSuite(t *testing.T) {
	suite.Run(t, &GrantCodeSuite{})
}

func (t *GrantCodeSuite) SetupTest() {
	t.store = NewMockallstore()
	t.grant = NewAuthorizationCode(t.store)
}

func (t *GrantCodeSuite) TestValidate() {
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

	t.store.On("GetAuthorizeCode", tr.Code).Return(ac, nil)

	grantData, resp := t.grant.Validate(tr)
	t.Nil(resp)
	t.EqualValues(grantData, expGrantData)
}

func (t *GrantCodeSuite) TestValidate_WithNotFound() {
	tr := &TokenRequest{
		Code:        "123",
		RedirectURI: "http://localhost",
	}

	t.store.On("GetAuthorizeCode", tr.Code).Return(nil, errors2.NewAnyError())

	grantData, resp := t.grant.Validate(tr)
	t.Equal(errAuthCodeNotExist, resp)
	t.Nil(grantData)
}

func (t *GrantCodeSuite) TestValidate_WithRedirectMisMatch() {
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

	t.store.On("GetAuthorizeCode", tr.Code).Return(ac, nil)

	grantData, resp := t.grant.Validate(tr)
	t.Equal(errRedirectMismatch, resp)
	t.Nil(grantData)
}

func (t *GrantCodeSuite) TestValidate_WithCodeExpired() {
	ac := &AuthorizeCode{
		ClientID:    "123",
		UserID:      "abc",
		Scope:       []string{"1", "2"},
		RedirectURI: "http://localhost",
		Expires:     addSecondUnix(0),
	}

	tr := &TokenRequest{
		Code:        "123",
		RedirectURI: "http://localhost",
	}

	t.store.On("GetAuthorizeCode", tr.Code).Return(ac, nil)

	grantData, resp := t.grant.Validate(tr)
	t.Equal(errAuthCodeExpired, resp)
	t.Nil(grantData)
}

func (t *GrantCodeSuite) TestValidate_WithCodeEmpty() {
	tr := &TokenRequest{
		Code:        "",
		RedirectURI: "http://localhost",
	}

	grantData, resp := t.grant.Validate(tr)
	t.Equal(errCodeRequired, resp)
	t.Nil(grantData)

}

func (t *GrantCodeSuite) TestName() {
	t.Equal(AUTHORIZATION_CODE, t.grant.Name())
}

func (t *GrantCodeSuite) TestIncludeRefreshToken() {
	t.True(t.grant.IncludeRefreshToken())
}

func (t *GrantCodeSuite) TestBeforeCreateAccessToken() {
	tr := &TokenRequest{}
	td := &TokenData{}
	t.Nil(t.grant.BeforeCreateAccessToken(tr, td))
}
