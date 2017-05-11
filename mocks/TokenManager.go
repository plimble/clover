package mocks

import clover "github.com/plimble/clover"
import mock "github.com/stretchr/testify/mock"

// TokenManager is an autogenerated mock type for the TokenManager type
type TokenManager struct {
	mock.Mock
}

// DeleteAccessToken provides a mock function with given fields: token
func (_m *TokenManager) DeleteAccessToken(token string) error {
	ret := _m.Called(token)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAuthorizeCode provides a mock function with given fields: code
func (_m *TokenManager) DeleteAuthorizeCode(code string) error {
	ret := _m.Called(code)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(code)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteRefreshToken provides a mock function with given fields: token
func (_m *TokenManager) DeleteRefreshToken(token string) error {
	ret := _m.Called(token)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateAccessToken provides a mock function with given fields: ctx, includeRefreshToken
func (_m *TokenManager) GenerateAccessToken(ctx *clover.AccessTokenContext, includeRefreshToken bool) (*clover.AccessToken, *clover.RefreshToken, error) {
	ret := _m.Called(ctx, includeRefreshToken)

	var r0 *clover.AccessToken
	if rf, ok := ret.Get(0).(func(*clover.AccessTokenContext, bool) *clover.AccessToken); ok {
		r0 = rf(ctx, includeRefreshToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clover.AccessToken)
		}
	}

	var r1 *clover.RefreshToken
	if rf, ok := ret.Get(1).(func(*clover.AccessTokenContext, bool) *clover.RefreshToken); ok {
		r1 = rf(ctx, includeRefreshToken)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*clover.RefreshToken)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*clover.AccessTokenContext, bool) error); ok {
		r2 = rf(ctx, includeRefreshToken)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GenerateAuthorizeCode provides a mock function with given fields: ctx, authorizeCodeLifespan
func (_m *TokenManager) GenerateAuthorizeCode(ctx *clover.AuthorizeContext, authorizeCodeLifespan int) (*clover.AuthorizeCode, error) {
	ret := _m.Called(ctx, authorizeCodeLifespan)

	var r0 *clover.AuthorizeCode
	if rf, ok := ret.Get(0).(func(*clover.AuthorizeContext, int) *clover.AuthorizeCode); ok {
		r0 = rf(ctx, authorizeCodeLifespan)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clover.AuthorizeCode)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*clover.AuthorizeContext, int) error); ok {
		r1 = rf(ctx, authorizeCodeLifespan)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateRefreshToken provides a mock function with given fields: ctx
func (_m *TokenManager) GenerateRefreshToken(ctx *clover.AccessTokenContext) (*clover.RefreshToken, error) {
	ret := _m.Called(ctx)

	var r0 *clover.RefreshToken
	if rf, ok := ret.Get(0).(func(*clover.AccessTokenContext) *clover.RefreshToken); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clover.RefreshToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*clover.AccessTokenContext) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccessToken provides a mock function with given fields: token
func (_m *TokenManager) GetAccessToken(token string) (*clover.AccessToken, error) {
	ret := _m.Called(token)

	var r0 *clover.AccessToken
	if rf, ok := ret.Get(0).(func(string) *clover.AccessToken); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clover.AccessToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAuthorizeCode provides a mock function with given fields: code
func (_m *TokenManager) GetAuthorizeCode(code string) (*clover.AuthorizeCode, error) {
	ret := _m.Called(code)

	var r0 *clover.AuthorizeCode
	if rf, ok := ret.Get(0).(func(string) *clover.AuthorizeCode); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clover.AuthorizeCode)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRefreshToken provides a mock function with given fields: token
func (_m *TokenManager) GetRefreshToken(token string) (*clover.RefreshToken, error) {
	ret := _m.Called(token)

	var r0 *clover.RefreshToken
	if rf, ok := ret.Get(0).(func(string) *clover.RefreshToken); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clover.RefreshToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
