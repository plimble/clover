package clover

import "github.com/stretchr/testify/mock"

type MockAccessTokenRespType struct {
	mock.Mock
}

func NewMockAccessTokenRespType() *MockAccessTokenRespType {
	return &MockAccessTokenRespType{}
}

func (m *MockAccessTokenRespType) Response(td *TokenData, includeRefresh bool) *Response {
	ret := m.Called(td, includeRefresh)

	var r0 *Response
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*Response)
	}

	return r0
}
