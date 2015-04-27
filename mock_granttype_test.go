package clover

import "github.com/stretchr/testify/mock"

type MockGrantType struct {
	mock.Mock
}

func NewMockGrantType() *MockGrantType {
	return &MockGrantType{}
}

func (m *MockGrantType) Validate(tr *TokenRequest) (*GrantData, *Response) {
	ret := m.Called(tr)

	var r0 *GrantData
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*GrantData)
	}
	var r1 *Response
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*Response)
	}

	return r0, r1
}
func (m *MockGrantType) Name() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *MockGrantType) IncludeRefreshToken() bool {
	ret := m.Called()

	r0 := ret.Get(0).(bool)

	return r0
}
func (m *MockGrantType) BeforeCreateAccessToken(tr *TokenRequest, td *TokenData) *Response {
	ret := m.Called(tr, td)

	var r0 *Response
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*Response)
	}

	return r0
}
