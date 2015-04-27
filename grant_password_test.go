package clover

import (
	"github.com/plimble/utils/errors2"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GrantPasswordSuite struct {
	suite.Suite
	store *Mockallstore
	grant GrantType
}

func TestGrantPasswordSuite(t *testing.T) {
	suite.Run(t, &GrantPasswordSuite{})
}

func (t *GrantPasswordSuite) TestValidate() {
	tr := &TokenRequest{
		Username: "test",
		Password: "xyz",
	}

	expGrantData := &GrantData{
		UserID: "001",
		Scope:  []string{"1", "2", "3"},
		Data:   map[string]interface{}{"a": 1, "b": "1"},
	}

	u := &TestUser{
		ID:       "001",
		Username: "test",
		Password: "xyz",
		Scope:    []string{"1", "2", "3"},
		Data:     map[string]interface{}{"a": 1, "b": "1"},
	}

	t.store.On("GetUser", tr.Username, tr.Password).Return(u, nil)

	grantData, resp := t.grant.Validate(tr)
	t.Nil(resp)
	t.EqualValues(expGrantData, grantData)
}

func (t *GrantPasswordSuite) TestValidate_WithNotFound() {
	tr := &TokenRequest{
		Username: "abc",
		Password: "xyz",
	}

	t.store.On("GetUser", tr.Username, tr.Password).Return(nil, errors2.NewAnyError())

	grantData, resp := t.grant.Validate(tr)
	t.Equal(errInvalidUsernamePAssword, resp)
	t.Nil(grantData)
}

func (t *GrantPasswordSuite) TestValidate_WithUsernamePasswordRequired() {
	tr := &TokenRequest{
		Username: "",
		Password: "",
	}

	grantData, resp := t.grant.Validate(tr)
	t.Equal(errUsernamePasswordRequired, resp)
	t.Nil(grantData)
}

func (t *GrantPasswordSuite) SetupTest() {
	t.store = NewMockallstore()
	t.grant = NewPassword(t.store)
}

func (t *GrantPasswordSuite) TestName() {
	t.Equal(PASSWORD, t.grant.Name())
}

func (t *GrantPasswordSuite) TestIncludeRefreshToken() {
	t.True(t.grant.IncludeRefreshToken())
}

func (t *GrantPasswordSuite) TestBeforeCreateAccessToken() {
	tr := &TokenRequest{}
	td := &TokenData{}
	t.Nil(t.grant.BeforeCreateAccessToken(tr, td))
}
