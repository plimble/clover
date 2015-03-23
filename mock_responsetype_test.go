package clover

import "github.com/stretchr/testify/mock"

type MockResponseType struct {
	mock.Mock
}

func NewMockResponseType() *MockResponseType {
	return &MockResponseType{}
}

func (m *MockResponseType) GetAuthResponse(ar *authorizeRequest, client Client, scopes []string) *Response {
	ret := m.Called(ar, client, scopes)

	var r0 *Response
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*Response)
	}

	return r0
}
func (m *MockResponseType) GetAccessToken(td *TokenData, includeRefresh bool) *Response {
	ret := m.Called(td, includeRefresh)

	var r0 *Response
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*Response)
	}

	return r0
}
