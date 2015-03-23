package clover

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockAuthCtrl struct {
	store *Mockallstore
}

func setupAuthCtrl() (*authController, *mockAuthCtrl) {
	store := NewMockallstore()
	config := NewAuthConfig(store)
	config.AddAuthCodeGrant(store)
	authRespType := newCodeRespType(config)
	tokenRespType := newTokenRespType(config)
	mock := &mockAuthCtrl{store}
	ctrl := newAuthController(config, authRespType, tokenRespType)

	return ctrl, mock
}

func TestValidateClientID(t *testing.T) {
	ctrl, m := setupAuthCtrl()

	ar := &authorizeRequest{clientID: "123"}

	m.store.On("GetClient", ar.clientID).Return(&DefaultClient{}, nil)

	resp := ctrl.validateClientID(ar)
	assert.Nil(t, resp)
}

func TestValidateClientIDWithEmptyClientID(t *testing.T) {
	ctrl, m := setupAuthCtrl()

	ar := &authorizeRequest{}

	m.store.On("GetClient", ar.clientID).Return(nil, nil)

	resp := ctrl.validateClientID(ar)
	assert.Equal(t, errNoClientID, resp)
}

func TestValidateClientIDWithInvalidClient(t *testing.T) {
	ctrl, m := setupAuthCtrl()

	ar := &authorizeRequest{clientID: "123"}

	m.store.On("GetClient", ar.clientID).Return(nil, errors.New("not found"))

	resp := ctrl.validateClientID(ar)
	assert.Equal(t, errInvalidClientID, resp)
}
