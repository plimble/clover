package memory

import (
	"errors"

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
	client        map[string]*oauth2.Client
	scope         map[string]*Scope
	accessToken   map[string]*oauth2.AccessToken
	refreshToken  map[string]*oauth2.RefreshToken
	authCodeToken map[string]*oauth2.AuthorizeCode
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) Flush() {
	s.client = make(map[string]*oauth2.Client)
	s.client = make(map[string]*oauth2.Client)
}

func (s *MemoryStorage) RevokeAccessToken(token string) error {
	_, ok := s.accessToken[token]
	if !ok {
		return errNotFound
	}

	delete(s.accessToken, token)

	return nil
}

func (s *MemoryStorage) SaveAccessToken(accessToken *oauth2.AccessToken) error {
	s.accessToken[accessToken.AccessToken] = accessToken
	return nil
}

func (s *MemoryStorage) GetAccessToken(token string) (*oauth2.AccessToken, error) {
	accessToken, ok := s.accessToken[token]
	if !ok {
		return nil, errNotFound
	}

	return accessToken, nil
}

func (s *MemoryStorage) RevokeRefreshToken(token string) error {
	_, ok := s.refreshToken[token]
	if !ok {
		return errNotFound
	}

	delete(s.refreshToken, token)

	return nil
}

func (s *MemoryStorage) SaveRefreshToken(refreshToken *oauth2.RefreshToken) error {
	s.refreshToken[refreshToken.RefreshToken] = refreshToken
	return nil
}

func (s *MemoryStorage) GetRefreshToken(token string) (*oauth2.RefreshToken, error) {
	refreshToken, ok := s.refreshToken[token]
	if !ok {
		return nil, errNotFound
	}

	return refreshToken, nil
}

func (s *MemoryStorage) RevokeAuthorizeCode(code string) error {
	_, ok := s.authCodeToken[code]
	if !ok {
		return errNotFound
	}

	delete(s.authCodeToken, code)

	return nil
}

func (s *MemoryStorage) SaveAuthorizeCode(authCode *oauth2.AuthorizeCode) error {
	s.authCodeToken[authCode.Code] = authCode
	return nil
}

func (s *MemoryStorage) GetAuthorizeCode(code string) (*oauth2.AuthorizeCode, error) {
	authCode, ok := s.authCodeToken[code]
	if !ok {
		return nil, errNotFound
	}

	return authCode, nil
}

func (s *MemoryStorage) DeleteClient(id string) error {
	_, ok := s.client[id]
	if !ok {
		return errNotFound
	}

	delete(s.client, id)

	return nil
}

func (s *MemoryStorage) SaveClient(client *oauth2.Client) error {
	s.client[client.ID] = client
	return nil
}

func (s *MemoryStorage) GetClientWithSecret(id, secret string) (*oauth2.Client, error) {
	client, ok := s.client[id]
	if !ok {
		return nil, errNotFound
	}

	if client.Secret != secret {
		return nil, errors.New("secret mismatch")
	}

	return client, nil
}

func (s *MemoryStorage) GetClient(id string) (*oauth2.Client, error) {
	client, ok := s.client[id]
	if !ok {
		return nil, errNotFound
	}

	return client, nil
}

func (s *MemoryStorage) CreateScope(scope *Scope) error {
	s.scope[scope.ID] = scope
	return nil
}

func (s *MemoryStorage) GetScopeByIDs(ids []string) ([]*Scope, error) {
	scopes := []*Scope{}

	for _, id := range ids {
		scope, ok := s.scope[id]
		if ok {
			scopes = append(scopes, scope)
		}
	}

	return scopes, nil
}

func (s *MemoryStorage) DeleteScope(id string) error {
	_, ok := s.scope[id]
	if !ok {
		return errNotFound
	}

	delete(s.scope, id)

	return nil
}

func (s *MemoryStorage) GetAllScope() ([]*Scope, error) {
	scopes := make([]*Scope, len(s.scope))
	index := 0
	for _, scope := range s.scope {
		scopes[index] = scope
		index++
	}

	return scopes, nil
}

func (s *MemoryStorage) IsAvailableScope(scopes []string) (bool, error) {
	for _, scope := range scopes {
		if _, ok := s.scope[scope]; !ok {
			return false, nil
		}
	}

	return true, nil
}
