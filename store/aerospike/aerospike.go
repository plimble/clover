package aerospike

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/plimble/aerosingle"
	. "github.com/plimble/clover"
	"github.com/plimble/utils/errors2"
)

var errNoPublicKey = errors2.NewInternal("No publickey in store")
var errNotFound = errors2.NewNotFound("not found")

type GetUserFunc func(username, password string) (string, []string, error)
type GetClientFunc func(clientID string) (Client, error)

type AeroStore struct {
	client        *aerosingle.Client
	key           *PublicKey
	getUserFunc   GetUserFunc
	getClientFunc GetClientFunc
}

func New(asClient *aerosingle.Client) *AeroStore {
	return &AeroStore{asClient, nil, nil, nil}
}

func (s *AeroStore) RegisterGetUserFunc(fn GetUserFunc) {
	s.getUserFunc = fn
}

func (s *AeroStore) RegisterGetClientFunc(fn GetClientFunc) {
	s.getClientFunc = fn
}

func (s *AeroStore) Close() {
	s.client.Close()
}

func (s *AeroStore) GetUser(username, password string) (string, []string, error) {
	if s.getUserFunc == nil {
		panic("Not implement GetUserFunc")
	}

	id, scopes, err := s.getUserFunc(username, password)
	return id, scopes, err
}

func (s *AeroStore) GetClient(cid string) (Client, error) {
	if s.getClientFunc == nil {
		panic("Not implement GetClientFunc")
	}

	client, err := s.getClientFunc(cid)
	return client, err
}

func (s *AeroStore) SetAccessToken(accessToken *AccessToken) error {
	policy := aerospike.NewWritePolicy(0, 0)
	policy.RecordExistsAction = aerospike.CREATE_ONLY

	return s.client.PutMsgPack(policy, "access_token", accessToken.AccessToken, accessToken)
}

func (s *AeroStore) GetAccessToken(at string) (*AccessToken, error) {
	var data AccessToken
	if err := s.client.GetMsgPack(nil, "access_token", at, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *AeroStore) SetRefreshToken(refreshToken *RefreshToken) error {
	policy := aerospike.NewWritePolicy(0, 0)
	policy.RecordExistsAction = aerospike.CREATE_ONLY

	return s.client.PutMsgPack(policy, "refresh_token", refreshToken.RefreshToken, refreshToken)
}

func (s *AeroStore) GetRefreshToken(rt string) (*RefreshToken, error) {
	var data RefreshToken
	if err := s.client.GetMsgPack(nil, "refresh_token", rt, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *AeroStore) RemoveRefreshToken(rt string) error {
	return s.client.Delete(nil, "refresh_token", rt)
}

func (s *AeroStore) SetAuthorizeCode(ac *AuthorizeCode) error {
	policy := aerospike.NewWritePolicy(0, 0)
	policy.RecordExistsAction = aerospike.CREATE_ONLY

	return s.client.PutMsgPack(policy, "auth_code", ac.Code, ac)
}

func (s *AeroStore) GetAuthorizeCode(code string) (*AuthorizeCode, error) {
	var data AuthorizeCode
	if err := s.client.GetMsgPack(nil, "auth_code", code, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *AeroStore) SetPublicKey(key *PublicKey) {
	s.key = key
}

func (s *AeroStore) GetKey(clientID string) (*PublicKey, error) {
	if s.key == nil {
		return nil, errNoPublicKey
	}

	return s.key, nil
}
