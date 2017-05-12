package clover_test

import (
	"errors"
	"testing"
	"time"

	"github.com/plimble/clover"
	"github.com/plimble/clover/mocks"
	"github.com/stretchr/testify/require"
)

func TestTokenManager(t *testing.T) {
	t.Run("GenerateAccessToken", testTokenGenerateAccessToken)
	t.Run("GenerateRefreshToken", testTokenGenerateRefreshToken)
	t.Run("GenerateAuthorizeCode", testTokenGenerateAuthorizeCode)
	t.Run("GetAccessToken", testTokenGetAccessToken)
	t.Run("GetRefreshToken", testTokenGetRefreshToken)
	t.Run("GetAuthorizeCode", testTokenGetAuthorizeCode)
	t.Run("DeleteAccessToken", testTokenDeleteAccessToken)
	t.Run("DeleteRefreshToken", testTokenDeleteRefreshToken)
	t.Run("DeleteAuthorizeCode", testTokenDeleteAuthorizeCode)
}

func addSecondUnix(sec int) int64 {
	return time.Now().UTC().Truncate(time.Nanosecond).Add(time.Second * time.Duration(sec)).Unix()
}

var accessTokenSampleData = []*clover.AccessToken{
	{"aaaa", "c1", "u1", addSecondUnix(3600), []string{"s1", "s2"}},
}

var refreshTokenSampleData = []*clover.RefreshToken{
	{"rrrr", "c1", "u1", addSecondUnix(4600), []string{"s1", "s2"}, 3600, 4600},
}

var authCodeSampleData = []*clover.AuthorizeCode{
	{"cccc", "c1", "u1", addSecondUnix(3600), []string{"s1", "s2"}, "http://example.com", "code"},
}

var accesstokenCtxSample = []*clover.AccessTokenContext{
	{
		Client:               clover.Client{ID: "c1"},
		UserID:               "u1",
		Scopes:               []string{"s1", "s2"},
		AccessTokenLifespan:  3600,
		RefreshTokenLifespan: 4600,
	},
}

var authcodeCtxSample = []*clover.AuthorizeContext{
	{
		Client:       clover.Client{ID: "c1"},
		UserID:       "u1",
		Scopes:       []string{"s1", "s2"},
		RedirectURI:  "http://example.com",
		ResponseType: "code",
	},
}

type tokenTest struct {
	storage *mocks.TokenStorage
	atGen   *mocks.TokenGenerator
	rfGen   *mocks.TokenGenerator
	acGen   *mocks.TokenGenerator
	manager clover.TokenManager
}

func setup() *tokenTest {
	m := &tokenTest{}
	m.storage = &mocks.TokenStorage{}
	m.atGen = &mocks.TokenGenerator{}
	m.rfGen = &mocks.TokenGenerator{}
	m.acGen = &mocks.TokenGenerator{}
	m.manager = clover.NewTokenManager(m.atGen, m.rfGen, m.acGen, m.storage)

	return m
}

func testTokenGenerateAccessToken(t *testing.T) {
	t.Run("Without RefreshToken", func(t *testing.T) {
		s := setup()
		s.atGen.On("Generate", "c1", "u1", "s1 s2", 3600).Return("aaaa", nil)
		s.storage.On("SaveAccessToken", accessTokenSampleData[0]).Return(nil)

		accessToken, refreshToken, err := s.manager.GenerateAccessToken(accesstokenCtxSample[0], false)
		require.NoError(t, err)
		require.Equal(t, accessTokenSampleData[0], accessToken)
		require.Nil(t, refreshToken)
		s.atGen.AssertExpectations(t)
		s.storage.AssertExpectations(t)
	})

	t.Run("With RefreshToken", func(t *testing.T) {
		s := setup()
		s.atGen.On("Generate", "c1", "u1", "s1 s2", 3600).Return("aaaa", nil)
		s.storage.On("SaveAccessToken", accessTokenSampleData[0]).Return(nil)
		s.rfGen.On("Generate", "c1", "u1", "s1 s2", 4600).Return("rrrr", nil)
		s.storage.On("SaveRefreshToken", refreshTokenSampleData[0]).Return(nil)

		accessToken, refreshToken, err := s.manager.GenerateAccessToken(accesstokenCtxSample[0], true)
		require.NoError(t, err)
		require.Equal(t, accessTokenSampleData[0], accessToken)
		require.Equal(t, refreshTokenSampleData[0], refreshToken)
		s.atGen.AssertExpectations(t)
		s.rfGen.AssertExpectations(t)
		s.storage.AssertExpectations(t)
	})
}

func testTokenGenerateRefreshToken(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		s := setup()
		s.rfGen.On("Generate", "c1", "u1", "s1 s2", 4600).Return("rrrr", nil)
		s.storage.On("SaveRefreshToken", refreshTokenSampleData[0]).Return(nil)

		refreshToken, err := s.manager.GenerateRefreshToken(accesstokenCtxSample[0])
		require.NoError(t, err)
		require.Equal(t, refreshTokenSampleData[0], refreshToken)
		s.rfGen.AssertExpectations(t)
		s.storage.AssertExpectations(t)
	})

	t.Run("Error on Generate", func(t *testing.T) {
		s := setup()
		expErr := errors.New("generate error")
		s.rfGen.On("Generate", "c1", "u1", "s1 s2", 4600).Return("", expErr)

		refreshToken, err := s.manager.GenerateRefreshToken(accesstokenCtxSample[0])
		require.Equal(t, expErr, err)
		require.Nil(t, refreshToken)
		s.rfGen.AssertExpectations(t)
	})

	t.Run("Error on save", func(t *testing.T) {
		s := setup()
		expErr := errors.New("save error")
		s.rfGen.On("Generate", "c1", "u1", "s1 s2", 4600).Return("rrrr", nil)
		s.storage.On("SaveRefreshToken", refreshTokenSampleData[0]).Return(expErr)

		refreshToken, err := s.manager.GenerateRefreshToken(accesstokenCtxSample[0])
		require.Equal(t, expErr, err)
		require.Equal(t, refreshTokenSampleData[0], refreshToken)
		s.rfGen.AssertExpectations(t)
		s.storage.AssertExpectations(t)
	})
}

func testTokenGenerateAuthorizeCode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		s := setup()

		s.acGen.On("Generate", "c1", "u1", "s1 s2", 3600).Return("cccc", nil)
		s.storage.On("SaveAuthorizeCode", authCodeSampleData[0]).Return(nil)

		authCode, err := s.manager.GenerateAuthorizeCode(authcodeCtxSample[0], 3600)
		require.NoError(t, err)
		require.Equal(t, authCodeSampleData[0], authCode)
		s.acGen.AssertExpectations(t)
		s.storage.AssertExpectations(t)
	})

	t.Run("Error on Generate", func(t *testing.T) {
		s := setup()

		expErr := errors.New("generate error")
		s.acGen.On("Generate", "c1", "u1", "s1 s2", 3600).Return("cccc", expErr)

		authCode, err := s.manager.GenerateAuthorizeCode(authcodeCtxSample[0], 3600)
		require.Equal(t, expErr, err)
		require.Nil(t, authCode)
		s.acGen.AssertExpectations(t)
	})
}

func testTokenGetAccessToken(t *testing.T) {
	s := setup()

	s.storage.On("GetAccessToken", "token").Return(accessTokenSampleData[0], nil)

	accessToken, err := s.manager.GetAccessToken("token")
	require.NoError(t, err)
	require.Equal(t, accessTokenSampleData[0], accessToken)
	s.storage.AssertExpectations(t)
}

func testTokenGetRefreshToken(t *testing.T) {
	s := setup()

	s.storage.On("GetRefreshToken", "token").Return(refreshTokenSampleData[0], nil)

	refreshToken, err := s.manager.GetRefreshToken("token")
	require.NoError(t, err)
	require.Equal(t, refreshTokenSampleData[0], refreshToken)
	s.storage.AssertExpectations(t)
}

func testTokenGetAuthorizeCode(t *testing.T) {
	s := setup()

	s.storage.On("GetAuthorizeCode", "code").Return(authCodeSampleData[0], nil)

	authCode, err := s.manager.GetAuthorizeCode("code")
	require.NoError(t, err)
	require.Equal(t, authCodeSampleData[0], authCode)
	s.storage.AssertExpectations(t)
}

func testTokenDeleteAccessToken(t *testing.T) {
	s := setup()

	s.storage.On("DeleteAccessToken", "code").Return(nil)

	err := s.manager.DeleteAccessToken("code")
	require.NoError(t, err)
	s.storage.AssertExpectations(t)
}

func testTokenDeleteRefreshToken(t *testing.T) {
	s := setup()

	s.storage.On("DeleteRefreshToken", "code").Return(nil)

	err := s.manager.DeleteRefreshToken("code")
	require.NoError(t, err)
	s.storage.AssertExpectations(t)
}

func testTokenDeleteAuthorizeCode(t *testing.T) {
	s := setup()

	s.storage.On("DeleteAuthorizeCode", "code").Return(nil)

	err := s.manager.DeleteAuthorizeCode("code")
	require.NoError(t, err)
	s.storage.AssertExpectations(t)
}
