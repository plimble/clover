package memory

import (
	"errors"
	"github.com/plimble/clover"
)

var (
	errNotFound = errors.New("not found")
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetUserFunc func(username, password string) (string, []string, error)
type GetClientFunc func(clientID string) (clover.Client, error)

type Storage struct {
	Client      map[string]*clover.DefaultClient
	Refresh     map[string]*clover.RefreshToken
	AuthCode    map[string]*clover.AuthorizeCode
	AccessToken map[string]*clover.AccessToken
	User        map[string]*User
}

func New() *Storage {
	return &Storage{
		Client:      make(map[string]*clover.DefaultClient),
		Refresh:     make(map[string]*clover.RefreshToken),
		AuthCode:    make(map[string]*clover.AuthorizeCode),
		AccessToken: make(map[string]*clover.AccessToken),
	}
}

func (s *Storage) GetClient(id string) (clover.Client, error) {
	client, ok := s.Client[id]
	if !ok {
		return nil, errNotFound
	}

	return client, nil
}

func (s *Storage) SetAccessToken(accessToken *clover.AccessToken) error {
	s.AccessToken[accessToken.AccessToken] = accessToken
	return nil
}

func (s *Storage) GetAccessToken(at string) (*clover.AccessToken, error) {
	accesstoken, ok := s.AccessToken[at]
	if !ok {
		return nil, errNotFound
	}

	return accesstoken, nil
}

func (s *Storage) SetRefreshToken(rt *clover.RefreshToken) error {
	s.Refresh[rt.RefreshToken] = rt
	return nil
}

func (s *Storage) GetRefreshToken(rt string) (*clover.RefreshToken, error) {
	refreshtoken, ok := s.Refresh[rt]
	if !ok {
		return nil, errNotFound
	}

	return refreshtoken, nil
}

func (s *Storage) RemoveRefreshToken(rt string) error {
	_, ok := s.Refresh[rt]
	if !ok {
		return errNotFound
	}

	delete(s.Refresh, rt)

	return nil
}

func (s *Storage) SetAuthorizeCode(ac *clover.AuthorizeCode) error {
	s.AuthCode[ac.Code] = ac
	return nil
}

func (s *Storage) GetAuthorizeCode(code string) (*clover.AuthorizeCode, error) {
	authcode, ok := s.AuthCode[code]
	if !ok {
		return nil, errNotFound
	}

	return authcode, nil
}

func (s *Storage) GetUser(username, password string) (string, []string, error) {
	user, ok := s.User[username]
	if !ok {
		return "", nil, errNotFound
	}

	if password != user.Password {
		return "", nil, errors.New("invalid username or password")
	}

	return user.Username, nil, nil
}
