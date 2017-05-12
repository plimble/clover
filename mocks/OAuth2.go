package mocks

import http "net/http"
import mock "github.com/stretchr/testify/mock"

// OAuth2 is an autogenerated mock type for the OAuth2 type
type OAuth2 struct {
	mock.Mock
}

// AccessTokenHandler provides a mock function with given fields: w, r
func (_m *OAuth2) AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// AuthorizeHandler provides a mock function with given fields: w, r
func (_m *OAuth2) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// IntrospectionHandler provides a mock function with given fields: w, r
func (_m *OAuth2) IntrospectionHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// RevokeHandler provides a mock function with given fields: w, r
func (_m *OAuth2) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}