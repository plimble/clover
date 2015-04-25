package aerospike

import (
	"github.com/plimble/aero"
	"github.com/plimble/clover"
	"github.com/plimble/utils/env"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StoreSuite struct {
	suite.Suite
	client *aero.Client
	store  *AeroStore
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, &StoreSuite{})
}

func (t *StoreSuite) SetupSuite() {
	t.client = aero.NewClient(env.String("AS_HOST", ""), 3000)
}

func (t *StoreSuite) SetupTest() {
	t.store = New(t.client, "test", 100, 100, 100)
}

func (t *StoreSuite) TearDownTest() {
	t.client.Delete(nil, t.store.ns, "access_token", "1")
	t.client.Delete(nil, t.store.ns, "auth_code", "1")
}

func (t *StoreSuite) TearDownSuite() {
	t.client.Close()
}

func (t *StoreSuite) TestAccessToken() {
	expa := &clover.AccessToken{
		AccessToken: "1",
		ClientID:    "2",
		UserID:      "3",
		Expires:     100,
		Scope:       []string{"1", "2"},
	}

	//put
	err := t.store.SetAccessToken(expa)
	t.NoError(err)

	//get
	a, err := t.store.GetAccessToken(expa.AccessToken)
	t.NoError(err)
	t.Equal(expa, a)

	//not found
	a, err = t.store.GetAccessToken("xxx")
	t.Equal(aero.ErrNotFound, err)
	t.Nil(a)
}

func (t *StoreSuite) TestRefreshToken() {
	expr := &clover.RefreshToken{
		RefreshToken: "1",
		ClientID:     "2",
		UserID:       "3",
		Expires:      100,
		Scope:        []string{"1", "2"},
	}

	//put
	err := t.store.SetRefreshToken(expr)
	t.NoError(err)

	//get
	r, err := t.store.GetRefreshToken(expr.RefreshToken)
	t.NoError(err)
	t.Equal(expr, r)

	//delete
	err = t.store.RemoveRefreshToken(expr.RefreshToken)
	t.NoError(err)

	//not found
	r, err = t.store.GetRefreshToken(expr.RefreshToken)
	t.Equal(aero.ErrNotFound, err)
	t.Nil(r)
}

func (t *StoreSuite) TestAuthorizeCode() {
	expa := &clover.AuthorizeCode{
		Code:        "1",
		ClientID:    "2",
		UserID:      "3",
		Expires:     100,
		Scope:       []string{"1", "2"},
		RedirectURI: "http://localhost",
	}

	//put
	err := t.store.SetAuthorizeCode(expa)
	t.NoError(err)

	//get
	a, err := t.store.GetAuthorizeCode(expa.Code)
	t.NoError(err)
	t.Equal(expa, a)
}

func (t *StoreSuite) TestGetKey() {
	expk := &clover.PublicKey{
		ClientID:   "123",
		PublicKey:  "abc",
		PrivateKey: "abc",
		Algorithm:  clover.JWT_ALGO_HS512,
	}

	//get without public key
	key, err := t.store.GetKey(expk.ClientID)
	t.Equal(errNoPublicKey, err)
	t.Nil(key)

	//get
	t.store.SetPublicKey(expk)
	key, err = t.store.GetKey(expk.ClientID)
	t.NoError(err)
	t.Equal(expk, key)
}
