package aerospike

import (
	"github.com/plimble/aero"
	. "github.com/plimble/clover"
	"github.com/plimble/utils/errors2"
)

var errNoPublicKey = errors2.NewInternal("No publickey in store")
var errNotFound = errors2.NewNotFound("not found")

type AeroStore struct {
	client           *aero.Client
	ns               string
	key              *PublicKey
	tokenLifeTime    int
	authCodeLifetime int
	refresLifeTime   int
}

func New(asClient *aero.Client, ns string, tokenLifeTime, authCodeLifetime, refresLifeTime int) *AeroStore {
	return &AeroStore{asClient, ns, nil, tokenLifeTime, authCodeLifetime, refresLifeTime}
}

func (s *AeroStore) SetAccessToken(accessToken *AccessToken) error {
	policy := aero.NewWritePolicy(0, int32(s.tokenLifeTime*2))
	policy.RecordExistsAction = aero.CREATE_ONLY

	b, err := aero.MarshalMsgPack(accessToken)
	if err != nil {
		return err
	}

	binAll := aero.NewBin("all", b)

	return s.client.PutBins(policy, s.ns, "access_token", accessToken.AccessToken, binAll)
}

func (s *AeroStore) GetAccessToken(at string) (*AccessToken, error) {
	rec, err := s.client.Get(nil, s.ns, "access_token", at, "all")
	if err != nil {
		return nil, err
	}

	token := &AccessToken{}
	if err := aero.UnmarshalMsgPack(rec.Bins["all"].([]byte), token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AeroStore) SetRefreshToken(refreshToken *RefreshToken) error {
	policy := aero.NewWritePolicy(0, int32(s.refresLifeTime*2))
	policy.RecordExistsAction = aero.CREATE_ONLY

	b, err := aero.MarshalMsgPack(refreshToken)
	if err != nil {
		return err
	}

	binAll := aero.NewBin("all", b)

	return s.client.PutBins(policy, s.ns, "refresh_token", refreshToken.RefreshToken, binAll)
}

func (s *AeroStore) GetRefreshToken(rt string) (*RefreshToken, error) {
	rec, err := s.client.Get(nil, s.ns, "refresh_token", rt, "all")
	if err != nil {
		return nil, err
	}

	token := &RefreshToken{}

	if err := aero.UnmarshalMsgPack(rec.Bins["all"].([]byte), token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AeroStore) RemoveRefreshToken(rt string) error {
	return s.client.Delete(nil, s.ns, "refresh_token", rt)
}

func (s *AeroStore) SetAuthorizeCode(ac *AuthorizeCode) error {
	policy := aero.NewWritePolicy(0, int32(s.authCodeLifetime*2))
	policy.RecordExistsAction = aero.CREATE_ONLY

	b, err := aero.MarshalMsgPack(ac)
	if err != nil {
		return err
	}

	binAll := aero.NewBin("all", b)

	return s.client.PutBins(policy, s.ns, "auth_code", ac.Code, binAll)
}

func (s *AeroStore) GetAuthorizeCode(code string) (*AuthorizeCode, error) {
	rec, err := s.client.Get(nil, s.ns, "auth_code", code, "all")
	if err != nil {
		return nil, err
	}

	authCode := &AuthorizeCode{}
	if err := aero.UnmarshalMsgPack(rec.Bins["all"].([]byte), authCode); err != nil {
		return nil, err
	}

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
