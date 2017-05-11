package memory

import (
	"sync"

	"github.com/plimble/clover"
)

type tokenStore struct {
	at            sync.Mutex
	rf            sync.Mutex
	ac            sync.Mutex
	accessToken   map[string]*clover.AccessToken
	refreshToken  map[string]*clover.RefreshToken
	authCodeToken map[string]*clover.AuthorizeCode
}

func NewTokenStorage() clover.TokenStorage {
	return &tokenStore{
		accessToken:   make(map[string]*clover.AccessToken),
		refreshToken:  make(map[string]*clover.RefreshToken),
		authCodeToken: make(map[string]*clover.AuthorizeCode),
	}
}

func (s *tokenStore) Flush() {
	s.accessToken = make(map[string]*clover.AccessToken)
	s.refreshToken = make(map[string]*clover.RefreshToken)
	s.authCodeToken = make(map[string]*clover.AuthorizeCode)
}

func (s *tokenStore) DeleteAccessToken(token string) error {
	s.at.Lock()
	defer s.at.Unlock()
	_, ok := s.accessToken[token]
	if !ok {
		return errNotFound
	}

	delete(s.accessToken, token)

	return nil
}

func (s *tokenStore) SaveAccessToken(accessToken *clover.AccessToken) error {
	s.at.Lock()
	defer s.at.Unlock()
	s.accessToken[accessToken.AccessToken] = accessToken
	return nil
}

func (s *tokenStore) GetAccessToken(token string) (*clover.AccessToken, error) {
	s.at.Lock()
	defer s.at.Unlock()
	accessToken, ok := s.accessToken[token]
	if !ok {
		return nil, errNotFound
	}

	return accessToken, nil
}

func (s *tokenStore) DeleteRefreshToken(token string) error {
	s.rf.Lock()
	defer s.rf.Unlock()
	_, ok := s.refreshToken[token]
	if !ok {
		return errNotFound
	}

	delete(s.refreshToken, token)

	return nil
}

func (s *tokenStore) SaveRefreshToken(refreshToken *clover.RefreshToken) error {
	s.rf.Lock()
	defer s.rf.Unlock()
	s.refreshToken[refreshToken.RefreshToken] = refreshToken
	return nil
}

func (s *tokenStore) GetRefreshToken(token string) (*clover.RefreshToken, error) {
	s.rf.Lock()
	defer s.rf.Unlock()
	refreshToken, ok := s.refreshToken[token]
	if !ok {
		return nil, errNotFound
	}

	return refreshToken, nil
}

func (s *tokenStore) DeleteAuthorizeCode(code string) error {
	s.ac.Lock()
	defer s.ac.Unlock()
	_, ok := s.authCodeToken[code]
	if !ok {
		return errNotFound
	}

	delete(s.authCodeToken, code)

	return nil
}

func (s *tokenStore) SaveAuthorizeCode(authCode *clover.AuthorizeCode) error {
	s.ac.Lock()
	defer s.ac.Unlock()
	s.authCodeToken[authCode.Code] = authCode
	return nil
}

func (s *tokenStore) GetAuthorizeCode(code string) (*clover.AuthorizeCode, error) {
	s.ac.Lock()
	defer s.ac.Unlock()
	authCode, ok := s.authCodeToken[code]
	if !ok {
		return nil, errNotFound
	}

	return authCode, nil
}
