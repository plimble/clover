package memory

import (
	"errors"

	"sync"

	"github.com/plimble/clover/oauth2"
)

type Scope struct {
	ID          string
	Description string
}

var (
	errNotFound = errors.New("not found")
)

type MemoryStorage struct {
	ClientMutex        sync.Mutex
	ScopeMutex         sync.Mutex
	AccessTokenMutex   sync.Mutex
	RefreshTokenMutex  sync.Mutex
	AuthorizeCodeMutex sync.Mutex
	Client             map[string]*oauth2.Client
	Scope              map[string]*Scope
	AccessToken        map[string]*oauth2.AccessToken
	RefreshToken       map[string]*oauth2.RefreshToken
	AuthorizeCode      map[string]*oauth2.AuthorizeCode
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Client:        make(map[string]*oauth2.Client),
		Scope:         make(map[string]*Scope),
		AccessToken:   make(map[string]*oauth2.AccessToken),
		RefreshToken:  make(map[string]*oauth2.RefreshToken),
		AuthorizeCode: make(map[string]*oauth2.AuthorizeCode),
	}
}

func (s *MemoryStorage) Flush() {
	s.Client = make(map[string]*oauth2.Client)
	s.Scope = make(map[string]*Scope)
	s.AccessToken = make(map[string]*oauth2.AccessToken)
	s.RefreshToken = make(map[string]*oauth2.RefreshToken)
	s.AuthorizeCode = make(map[string]*oauth2.AuthorizeCode)
}

func (s *MemoryStorage) RevokeAccessToken(token string) error {
	s.AccessTokenMutex.Lock()
	defer s.AccessTokenMutex.Unlock()

	_, ok := s.AccessToken[token]
	if !ok {
		return errNotFound
	}

	delete(s.AccessToken, token)

	return nil
}

func (s *MemoryStorage) SaveAccessToken(accessToken *oauth2.AccessToken) error {
	s.AccessTokenMutex.Lock()
	defer s.AccessTokenMutex.Unlock()

	s.AccessToken[accessToken.AccessToken] = accessToken
	return nil
}

func (s *MemoryStorage) GetAccessToken(token string) (*oauth2.AccessToken, error) {
	s.AccessTokenMutex.Lock()
	defer s.AccessTokenMutex.Unlock()

	accessToken, ok := s.AccessToken[token]
	if !ok {
		return nil, errNotFound
	}

	return accessToken, nil
}

func (s *MemoryStorage) RevokeRefreshToken(token string) error {
	s.RefreshTokenMutex.Lock()
	defer s.RefreshTokenMutex.Unlock()

	_, ok := s.RefreshToken[token]
	if !ok {
		return errNotFound
	}

	delete(s.RefreshToken, token)

	return nil
}

func (s *MemoryStorage) SaveRefreshToken(refreshToken *oauth2.RefreshToken) error {
	s.RefreshTokenMutex.Lock()
	defer s.RefreshTokenMutex.Unlock()

	s.RefreshToken[refreshToken.RefreshToken] = refreshToken
	return nil
}

func (s *MemoryStorage) GetRefreshToken(token string) (*oauth2.RefreshToken, error) {
	s.RefreshTokenMutex.Lock()
	defer s.RefreshTokenMutex.Unlock()

	refreshToken, ok := s.RefreshToken[token]
	if !ok {
		return nil, errNotFound
	}

	return refreshToken, nil
}

func (s *MemoryStorage) SaveAuthorizeCode(authCode *oauth2.AuthorizeCode) error {
	s.AuthorizeCodeMutex.Lock()
	defer s.AuthorizeCodeMutex.Unlock()

	s.AuthorizeCode[authCode.Code] = authCode
	return nil
}

func (s *MemoryStorage) GetAuthorizeCode(code string) (*oauth2.AuthorizeCode, error) {
	s.AuthorizeCodeMutex.Lock()
	defer s.AuthorizeCodeMutex.Unlock()

	authCode, ok := s.AuthorizeCode[code]
	if !ok {
		return nil, errNotFound
	}

	return authCode, nil
}

func (s *MemoryStorage) GetClientWithSecret(id, secret string) (*oauth2.Client, error) {
	s.ClientMutex.Lock()
	defer s.ClientMutex.Unlock()

	client, ok := s.Client[id]
	if !ok {
		return nil, errNotFound
	}

	if client.Secret != secret {
		return nil, errors.New("secret mismatch")
	}

	return client, nil
}

func (s *MemoryStorage) GetClient(id string) (*oauth2.Client, error) {
	s.ClientMutex.Lock()
	defer s.ClientMutex.Unlock()

	client, ok := s.Client[id]
	if !ok {
		return nil, errNotFound
	}

	return client, nil
}

func (s *MemoryStorage) IsAvailableScope(scopes []string) (bool, error) {
	s.ScopeMutex.Lock()
	defer s.ScopeMutex.Unlock()

	for _, scope := range scopes {
		if _, ok := s.Scope[scope]; !ok {
			return false, nil
		}
	}

	return true, nil
}
