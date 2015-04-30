package clover

import "github.com/stretchr/testify/mock"

type Mockallstore struct {
	mock.Mock
}

func NewMockallstore() *Mockallstore {
	return &Mockallstore{}
}

func (m *Mockallstore) GetClient(id string) (Client, error) {
	ret := m.Called(id)

	var r0 Client
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(Client)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Mockallstore) GetUser(username string, password string) (User, error) {
	ret := m.Called(username, password)

	var r0 User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(User)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Mockallstore) SetAccessToken(accessToken *AccessToken) error {
	ret := m.Called(accessToken)

	r0 := ret.Error(0)

	return r0
}
func (m *Mockallstore) GetAccessToken(at string) (*AccessToken, error) {
	ret := m.Called(at)

	var r0 *AccessToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*AccessToken)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Mockallstore) RemoveRefreshToken(rt string) error {
	ret := m.Called(rt)

	r0 := ret.Error(0)

	return r0
}
func (m *Mockallstore) SetRefreshToken(rt *RefreshToken) error {
	ret := m.Called(rt)

	r0 := ret.Error(0)

	return r0
}
func (m *Mockallstore) GetRefreshToken(rt string) (*RefreshToken, error) {
	ret := m.Called(rt)

	var r0 *RefreshToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*RefreshToken)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Mockallstore) SetAuthorizeCode(ac *AuthorizeCode) error {
	ret := m.Called(ac)

	r0 := ret.Error(0)

	return r0
}
func (m *Mockallstore) GetAuthorizeCode(code string) (*AuthorizeCode, error) {
	ret := m.Called(code)

	var r0 *AuthorizeCode
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*AuthorizeCode)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Mockallstore) GetKey(clientID string) (*PublicKey, error) {
	ret := m.Called(clientID)

	var r0 *PublicKey
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*PublicKey)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Mockallstore) ExistScopes(scopes ...string) (bool, error) {
	ret := m.Called(scopes)

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Mockallstore) GetDefaultScope(clientID string) ([]string, error) {
	ret := m.Called(clientID)

	var r0 []string
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]string)
	}
	r1 := ret.Error(1)

	return r0, r1
}
