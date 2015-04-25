package clover

import (
	"errors"
	"github.com/plimble/unik/mock_unik"
	"github.com/stretchr/testify/assert"
	testmock "github.com/stretchr/testify/mock"
	"testing"
)

type mockJWTRespType struct {
	store *Mockallstore
	unik  *mock_unik.MockGenerator
}

func setupJWTRespType() (*jwtResponseType, *mockJWTRespType) {
	store := NewMockallstore()
	unik := mock_unik.NewMockGenerator()

	config := NewAuthConfig(store)
	config.UseJWTAccessTokens(store)
	config.AuthCodeStore = store

	a := &AuthServer{
		config: config,
	}
	tokenRespType := newTokenRespType(a.config, unik)

	mock := &mockJWTRespType{store, unik}
	rt := newJWTResponseType(config, unik, tokenRespType)
	return rt, mock
}

func TestJWTRespType_CreateToken_WithRSKEY(t *testing.T) {
	rt, mock := setupJWTRespType()

	key := &PublicKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Algorithm:  JWT_ALGO_RS512,
	}

	mock.unik.On("Generate").Return("1")
	mock.store.On("GetKey", "1001").Return(key, nil)

	mock.store.On("SetAccessToken", testmock.Anything).Return(nil)
	at, resp := rt.createToken("1001", "1", []string{"email"}, nil)

	assert.Nil(t, resp)
	assert.NotEmpty(t, at.AccessToken)
}

func TestJWTRespType_CreateToken_WithHSKEY(t *testing.T) {
	rt, mock := setupJWTRespType()

	key := &PublicKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Algorithm:  JWT_ALGO_HS512,
	}

	mock.unik.On("Generate").Return("1")
	mock.store.On("GetKey", "1001").Return(key, nil)

	mock.store.On("SetAccessToken", testmock.Anything).Return(nil)
	at, resp := rt.createToken("1001", "1", []string{"email"}, nil)

	assert.Nil(t, resp)
	assert.NotEmpty(t, at.AccessToken)
}

func TestJWTRespType_CreateToken_WithRSKEY_ErrorGetKey(t *testing.T) {
	rt, mock := setupJWTRespType()

	key := &PublicKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Algorithm:  JWT_ALGO_RS512,
	}

	mock.unik.On("Generate").Return("1")
	mock.store.On("GetKey", "1001").Return(key, errors.New("test"))

	mock.store.On("SetAccessToken", testmock.Anything).Return(nil)
	at, resp := rt.createToken("1001", "1", []string{"email"}, nil)

	assert.Nil(t, at)
	assert.Equal(t, 500, resp.code)
}

func TestJWTRespType_CreateToken_WithRSKEY_ErrorSetAccessToken(t *testing.T) {
	rt, mock := setupJWTRespType()

	key := &PublicKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Algorithm:  JWT_ALGO_RS512,
	}

	mock.unik.On("Generate").Return("1")
	mock.store.On("GetKey", "1001").Return(key, nil)

	mock.store.On("SetAccessToken", testmock.Anything).Return(errors.New("test"))
	at, resp := rt.createToken("1001", "1", []string{"email"}, nil)

	assert.Nil(t, at)
	assert.Equal(t, 500, resp.code)
}
