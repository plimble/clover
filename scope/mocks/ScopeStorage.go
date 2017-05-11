package mocks

import clover "github.com/plimble/clover"
import mock "github.com/stretchr/testify/mock"

// ScopeStorage is an autogenerated mock type for the ScopeStorage type
type ScopeStorage struct {
	mock.Mock
}

// GetScopeByIDs provides a mock function with given fields: ids
func (_m *ScopeStorage) GetScopeByIDs(ids []string) ([]*clover.Scope, error) {
	ret := _m.Called(ids)

	var r0 []*clover.Scope
	if rf, ok := ret.Get(0).(func([]string) []*clover.Scope); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*clover.Scope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
