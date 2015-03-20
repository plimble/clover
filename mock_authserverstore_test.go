package clover

import "github.com/stretchr/testify/mock"

type MockAuthServerStore struct {
	mock.Mock
}

func NewMockAuthServerStore() *MockAuthServerStore {
	return &MockAuthServerStore{}
}

func (m *MockAuthServerStore) GetClient(id string) (Client, error) {
	ret := m.Called(id)

	r0 := ret.Get(0).(Client)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *MockAuthServerStore) RemoveRefreshToken(rt string) error {
	ret := m.Called(rt)

	r0 := ret.Error(0)

	return r0
}
func (m *MockAuthServerStore) GetUser(username string, password string) (string, []string, error) {
	ret := m.Called(username, password)

	r0 := ret.Get(0).(string)
	r1 := ret.Get(1).([]string)
	r2 := ret.Error(2)

	return r0, r1, r2
}
func (m *MockAuthServerStore) SetAuthorizeCode(ac *AuthorizeCode) error {
	ret := m.Called(ac)

	r0 := ret.Error(0)

	return r0
}
func (m *MockAuthServerStore) GetAuthorizeCode(code string) (*AuthorizeCode, error) {
	ret := m.Called(code)

	r0 := ret.Get(0).(*AuthorizeCode)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *MockAuthServerStore) SetRefreshToken(rt *RefreshToken) error {
	ret := m.Called(rt)

	r0 := ret.Error(0)

	return r0
}
func (m *MockAuthServerStore) GetRefreshToken(rt string) (*RefreshToken, error) {
	ret := m.Called(rt)

	r0 := ret.Get(0).(*RefreshToken)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *MockAuthServerStore) SetAccessToken(accessToken *AccessToken) error {
	ret := m.Called(accessToken)

	r0 := ret.Error(0)

	return r0
}
func (m *MockAuthServerStore) GetAccessToken(at string) (*AccessToken, error) {
	ret := m.Called(at)

	r0 := ret.Get(0).(*AccessToken)
	r1 := ret.Error(1)

	return r0, r1
}
