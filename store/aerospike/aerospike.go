package aerospike

import (
	"github.com/aerospike/aerospike-client-go"
	. "github.com/plimble/clover"
	"github.com/plimble/utils/assingle"
	"github.com/plimble/utils/errors2"
)

var errNoPublicKey = errors2.NewInternal("No publickey in store")
var errNotFound = errors2.NewNotFound("not found")

type GetUserFunc func(username, password string) (string, []string, error)
type GetClientFunc func(clientID string) (Client, error)

type Store struct {
	single        *assingle.ASSingle
	key           *PublicKey
	getUserFunc   GetUserFunc
	getClientFunc GetClientFunc
}

func New(client *aerospike.Client, ns string) *Store {
	return &Store{assingle.New(client, ns), nil, nil, nil}
}

func (s *Store) RegisterGetUserFunc(fn GetUserFunc) {
	s.getUserFunc = fn
}

func (s *Store) RegisterGetClientFunc(fn GetClientFunc) {
	s.getClientFunc = fn
}

func (s *Store) Close() {
	s.single.Close()
}

func (s *Store) GetUser(username, password string) (string, []string, error) {
	if s.getUserFunc == nil {
		panic("Not implement GetUserFunc")
	}

	id, scopes, err := s.getUserFunc(username, password)
	return id, scopes, err
}

func (s *Store) GetClient(cid string) (Client, error) {
	if s.getClientFunc == nil {
		panic("Not implement GetClientFunc")
	}

	client, err := s.getClientFunc(cid)
	return client, err
}

func (s *Store) SetAccessToken(accessToken *AccessToken) error {
	policy := aerospike.NewWritePolicy(0, 0)
	policy.RecordExistsAction = aerospike.CREATE_ONLY

	return s.single.Put(policy, "access_token", accessToken.AccessToken, accessToken)
}

func (s *Store) GetAccessToken(at string) (*AccessToken, error) {
	var data AccessToken
	if err := s.single.Get(nil, "access_token", at, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *Store) SetRefreshToken(refreshToken *RefreshToken) error {
	policy := aerospike.NewWritePolicy(0, 0)
	policy.RecordExistsAction = aerospike.CREATE_ONLY

	return s.single.Put(policy, "refresh_token", refreshToken.RefreshToken, refreshToken)
}

func (s *Store) GetRefreshToken(rt string) (*RefreshToken, error) {
	var data RefreshToken
	if err := s.single.Get(nil, "refresh_token", rt, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *Store) RemoveRefreshToken(rt string) error {
	return s.single.Delete(nil, "refresh_token", rt)
}

func (s *Store) SetAuthorizeCode(ac *AuthorizeCode) error {
	policy := aerospike.NewWritePolicy(0, 0)
	policy.RecordExistsAction = aerospike.CREATE_ONLY

	return s.single.Put(policy, "auth_code", ac.Code, ac)
}

func (s *Store) GetAuthorizeCode(code string) (*AuthorizeCode, error) {
	var data AuthorizeCode
	if err := s.single.Get(nil, "auth_code", code, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *Store) SetPublicKey(key *PublicKey) {
	s.key = key
}

func (s *Store) GetKey(clientID string) (*PublicKey, error) {
	if s.key == nil {
		return nil, errNoPublicKey
	}

	return s.key, nil
}
