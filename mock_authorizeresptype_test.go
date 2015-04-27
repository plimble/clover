package clover

import "github.com/stretchr/testify/mock"

type MockAuthorizeRespType struct {
	mock.Mock
}

func NewMockAuthorizeRespType() *MockAuthorizeRespType {
	return &MockAuthorizeRespType{}
}

func (m *MockAuthorizeRespType) Name() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *MockAuthorizeRespType) Response(ad *AuthorizeData) *Response {
	ret := m.Called(ad)

	var r0 *Response
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*Response)
	}

	return r0
}
func (m *MockAuthorizeRespType) SupportGrant() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *MockAuthorizeRespType) IsImplicit() bool {
	ret := m.Called()

	r0 := ret.Get(0).(bool)

	return r0
}
