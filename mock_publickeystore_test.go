package clover

import "github.com/stretchr/testify/mock"

type MockPublicKeyStore struct {
	mock.Mock
}

func NewMockPublicKeyStore() *MockPublicKeyStore {
	return &MockPublicKeyStore{}
}

func (m *MockPublicKeyStore) GetKey(clientID string) (*PublicKey, error) {
	ret := m.Called(clientID)

	r0 := ret.Get(0).(*PublicKey)
	r1 := ret.Error(1)

	return r0, r1
}
