package memory

import (
	"errors"
	. "github.com/plimble/clover"
	"io/ioutil"
)

var (
	errNotFound = errors.New("not found")
)

type Store struct {
	User         map[string]User
	Client       map[string]Client
	Refresh      map[string]*RefreshToken
	AuthCode     map[string]*AuthorizeCode
	AccessToken  map[string]*AccessToken
	Scope        map[string]struct{}
	defaultScope []string
	PublicKey    map[string]*PublicKey
	cPublicKey   string
	cPrivateKey  string
	cHmacKey     string
}

func New() *Store {
	return &Store{
		User:        make(map[string]User),
		Client:      make(map[string]Client),
		Refresh:     make(map[string]*RefreshToken),
		AuthCode:    make(map[string]*AuthorizeCode),
		AccessToken: make(map[string]*AccessToken),
		PublicKey:   make(map[string]*PublicKey),
	}
}

func (s *Store) Flush() {
	s.User = make(map[string]User)
	s.Client = make(map[string]Client)
	s.Refresh = make(map[string]*RefreshToken)
	s.AuthCode = make(map[string]*AuthorizeCode)
	s.AccessToken = make(map[string]*AccessToken)
	s.PublicKey = make(map[string]*PublicKey)
	s.Scope = make(map[string]struct{})
}

func (s *Store) SetClient(c Client) error {
	s.Client[c.GetClientID()] = c
	return nil
}

func (s *Store) GetClient(id string) (Client, error) {
	client, ok := s.Client[id]
	if !ok {
		return nil, errNotFound
	}

	return client, nil
}

func (s *Store) SetAccessToken(accessToken *AccessToken) error {
	s.AccessToken[accessToken.AccessToken] = accessToken
	return nil
}

func (s *Store) GetAccessToken(at string) (*AccessToken, error) {
	accesstoken, ok := s.AccessToken[at]
	if !ok {
		return nil, errNotFound
	}

	return accesstoken, nil
}

func (s *Store) SetRefreshToken(rt *RefreshToken) error {
	s.Refresh[rt.RefreshToken] = rt
	return nil
}

func (s *Store) GetRefreshToken(rt string) (*RefreshToken, error) {
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

func (s *Store) SetAuthorizeCode(ac *AuthorizeCode) error {
	s.AuthCode[ac.Code] = ac
	return nil
}

func (s *Store) GetAuthorizeCode(code string) (*AuthorizeCode, error) {
	authcode, ok := s.AuthCode[code]
	if !ok {
		return nil, errNotFound
	}

	return authcode, nil
}

func (s *Store) SetUser(u User) error {
	s.User[u.GetUsername()] = u
	return nil
}

func (s *Store) GetUser(username, password string) (User, error) {
	user, ok := s.User[username]
	if !ok {
		return nil, errNotFound
	}

	if password != user.GetPassword() {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *Store) GetKey(clientID string) (*PublicKey, error) {
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
		s.PublicKey[clientID] = &PublicKey{
			PublicKey:  s.cPublicKey,
			PrivateKey: s.cPrivateKey,
			Algorithm:  JWT_ALGO_RS512,
		}
	}
}

func (s *Store) AddHSKey(clientID string) {
	if s.cHmacKey == "" {
		s.setHmacKey()
	}

	c, _ := s.GetClient(clientID)
	if c != nil {
		s.PublicKey[clientID] = &PublicKey{
			PublicKey:  s.cHmacKey,
			PrivateKey: s.cHmacKey,
			Algorithm:  JWT_ALGO_HS512,
		}
	}
}

func (s *Store) SetScope(key string) {
	s.Scope[key] = struct{}{}
}

func (s *Store) ExistScopes(scopes ...string) (bool, error) {
	for i := 0; i < len(scopes); i++ {
		if _, ok := s.Scope[scopes[i]]; !ok {
			return false, nil
		}
	}

	return true, nil
}

func (s *Store) SetDefaultScope(scope ...string) {
	s.defaultScope = scope
}

func (s *Store) GetDefaultScope(clientID string) ([]string, error) {
	return s.defaultScope, nil
}
