package aerospike

import (
	"github.com/aerospike/aerospike-client-go"
	. "github.com/plimble/clover"
	"github.com/plimble/utils/ashelper"
	"github.com/plimble/utils/errors2"
)

var errNoPublicKey = errors2.NewInternal("No publickey in store")
var errNotFound = errors2.NewNotFound("not found")

type GetUserFunc func(username, password string) (string, []string, error)
type GetClientFunc func(clientID string) (Client, error)

type AeroStore struct {
	client        *aerospike.Client
	ns            string
	key           *PublicKey
	getUserFunc   GetUserFunc
	getClientFunc GetClientFunc
}

func New(asClient *aerospike.Client, ns string) *AeroStore {
	return &AeroStore{asClient, ns, nil, nil, nil}
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

	key, _ := aerospike.NewKey(s.ns, "access_token", accessToken.AccessToken)

	binAll := aerospike.NewBin("all", ashelper.MarshalMsgPack(accessToken))

	return ashelper.ErrPut(s.client.PutBins(policy, key, binAll))
}

func (s *AeroStore) GetAccessToken(at string) (*AccessToken, error) {
	key, _ := aerospike.NewKey(s.ns, "access_token", at)

	rec, err := s.client.Get(nil, key)
	if err := ashelper.ErrGet(rec, err); err != nil {
		return nil, err
	}

	token := &AccessToken{}
	ashelper.UnmarshalMsgPack(rec.Bins["all"].([]byte), token)

	return token, nil
}

func (s *AeroStore) SetRefreshToken(refreshToken *RefreshToken) error {
	policy := aerospike.NewWritePolicy(0, 0)
	policy.RecordExistsAction = aerospike.CREATE_ONLY

	key, _ := aerospike.NewKey(s.ns, "refresh_token", refreshToken.RefreshToken)

	binAll := aerospike.NewBin("all", ashelper.MarshalMsgPack(refreshToken))

	return ashelper.ErrPut(s.client.PutBins(policy, key, binAll))
}

func (s *AeroStore) GetRefreshToken(rt string) (*RefreshToken, error) {
	key, _ := aerospike.NewKey(s.ns, "refresh_token", rt)

	rec, err := s.client.Get(nil, key)
	if err := ashelper.ErrGet(rec, err); err != nil {
		return nil, err
	}

	token := &RefreshToken{}
	ashelper.UnmarshalMsgPack(rec.Bins["all"].([]byte), token)

	return token, nil
}

func (s *AeroStore) RemoveRefreshToken(rt string) error {
	key, _ := aerospike.NewKey(s.ns, "refresh_token", rt)
	return ashelper.ErrDel(s.client.Delete(nil, key))
}

func (s *AeroStore) SetAuthorizeCode(ac *AuthorizeCode) error {
	policy := aerospike.NewWritePolicy(0, 0)
	policy.RecordExistsAction = aerospike.CREATE_ONLY

	key, _ := aerospike.NewKey(s.ns, "auth_code", ac.Code)

	binAll := aerospike.NewBin("all", ashelper.MarshalMsgPack(ac))

	return ashelper.ErrPut(s.client.PutBins(policy, key, binAll))
}

func (s *AeroStore) GetAuthorizeCode(code string) (*AuthorizeCode, error) {
	key, _ := aerospike.NewKey(s.ns, "auth_code", code)

	rec, err := s.client.Get(nil, key)
	if err := ashelper.ErrGet(rec, err); err != nil {
		return nil, err
	}

	authCode := &AuthorizeCode{}
	ashelper.UnmarshalMsgPack(rec.Bins["all"].([]byte), authCode)

	return authCode, nil
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
