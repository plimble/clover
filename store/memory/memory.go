package memory

import (
	"errors"
	"github.com/plimble/clover"
	"io/ioutil"
)

var (
	errNotFound = errors.New("not found")
)

type Store struct {
	User        map[string]clover.User
	Client      map[string]clover.Client
	Refresh     map[string]*clover.RefreshToken
	AuthCode    map[string]*clover.AuthorizeCode
	AccessToken map[string]*clover.AccessToken
	PublicKey   map[string]*clover.PublicKey
	cPublicKey  string
	cPrivateKey string
	cHmacKey    string
}

func New() *Store {
	return &Store{
		User:        make(map[string]clover.User),
		Client:      make(map[string]clover.Client),
		Refresh:     make(map[string]*clover.RefreshToken),
		AuthCode:    make(map[string]*clover.AuthorizeCode),
		AccessToken: make(map[string]*clover.AccessToken),
		PublicKey:   make(map[string]*clover.PublicKey),
	}
}

func (s *Store) Flush() {
	s.User = make(map[string]clover.User)
	s.Client = make(map[string]clover.Client)
	s.Refresh = make(map[string]*clover.RefreshToken)
	s.AuthCode = make(map[string]*clover.AuthorizeCode)
	s.AccessToken = make(map[string]*clover.AccessToken)
	s.PublicKey = make(map[string]*clover.PublicKey)
}

func (s *Store) SetClient(c clover.Client) error {
	s.Client[c.GetClientID()] = c
	return nil
}

func (s *Store) GetClient(id string) (clover.Client, error) {
	client, ok := s.Client[id]
	if !ok {
		return nil, errNotFound
	}

	return client, nil
}

func (s *Store) SetAccessToken(accessToken *clover.AccessToken) error {
	s.AccessToken[accessToken.AccessToken] = accessToken
	return nil
}

func (s *Store) GetAccessToken(at string) (*clover.AccessToken, error) {
	accesstoken, ok := s.AccessToken[at]
	if !ok {
		return nil, errNotFound
	}

	return accesstoken, nil
}

func (s *Store) SetRefreshToken(rt *clover.RefreshToken) error {
	s.Refresh[rt.RefreshToken] = rt
	return nil
}

func (s *Store) GetRefreshToken(rt string) (*clover.RefreshToken, error) {
	refreshtoken, ok := s.Refresh[rt]
	if !ok {
		return nil, errNotFound
	}

	return refreshtoken, nil
}

func (s *Store) RemoveRefreshToken(rt string) error {
	_, ok := s.Refresh[rt]
	if !ok {
		return errNotFound
	}

	delete(s.Refresh, rt)

	return nil
}

func (s *Store) SetAuthorizeCode(ac *clover.AuthorizeCode) error {
	s.AuthCode[ac.Code] = ac
	return nil
}

func (s *Store) GetAuthorizeCode(code string) (*clover.AuthorizeCode, error) {
	authcode, ok := s.AuthCode[code]
	if !ok {
		return nil, errNotFound
	}

	return authcode, nil
}

func (s *Store) SetUser(u clover.User) error {
	s.User[u.GetUsername()] = u
	return nil
}

func (s *Store) GetUser(username, password string) (clover.User, error) {
	user, ok := s.User[username]
	if !ok {
		return nil, errNotFound
	}

	if password != user.GetPassword() {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *Store) GetKey(clientID string) (*clover.PublicKey, error) {
	pub, ok := s.PublicKey[clientID]
	if !ok {
		return nil, errNotFound
	}

	return pub, nil
}

func (s *Store) setHmacKey() {
	b, _ := ioutil.ReadFile("../tests/jwt_test/hmac")
	s.cHmacKey = string(b)
}

func (s *Store) setRSKey() {
	b, _ := ioutil.ReadFile("../tests/jwt_test/sample_key")
	s.cPrivateKey = string(b)

	b, _ = ioutil.ReadFile("../tests/jwt_test/key.pub")
	s.cPublicKey = string(b)
}

func (s *Store) AddRSKey(clientID string) {
	if s.cPublicKey == "" || s.cPrivateKey == "" {
		s.setRSKey()
	}

	c, _ := s.GetClient(clientID)
	if c != nil {
		s.PublicKey[clientID] = &clover.PublicKey{
			PublicKey:  s.cPublicKey,
			PrivateKey: s.cPrivateKey,
			Algorithm:  clover.JWT_ALGO_RS512,
		}
	}
}

func (s *Store) AddHSKey(clientID string) {
	if s.cHmacKey == "" {
		s.setHmacKey()
	}

	c, _ := s.GetClient(clientID)
	if c != nil {
		s.PublicKey[clientID] = &clover.PublicKey{
			PublicKey:  s.cHmacKey,
			PrivateKey: s.cHmacKey,
			Algorithm:  clover.JWT_ALGO_HS512,
		}
	}
}
