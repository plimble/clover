package mongo

import (
	"github.com/plimble/clover"
	"github.com/plimble/utils/errors2/errmgo"
	"gopkg.in/mgo.v2"
)

type Storage struct {
	session       *mgo.Session
	db            string
	getUserFunc   GetUserFunc
	getClientFunc GetClientFunc
}

type GetUserFunc func(username, password string) (string, []string, error)
type GetClientFunc func(clientID string) (clover.Client, error)

func New(session *mgo.Session, db string) *Storage {
	return &Storage{
		session: session,
		db:      db,
	}
}

func (s *Storage) RegisterGetUserFunc(fn GetUserFunc) {
	s.getUserFunc = fn
}

func (s *Storage) RegisterGetClientFunc(fn GetClientFunc) {
	s.getClientFunc = fn
}

func (s *Storage) TruncateAll() {
	session := s.session.Copy()
	defer session.Close()
	// session.DB(s.db).C("activity").RemoveAll(nil)
}

func (s *Storage) GetUser(username, password string) (string, []string, error) {
	session := s.session.Copy()
	defer session.Close()

	if s.getUserFunc == nil {
		panic("Not implement GetUserFunc")
	}

	id, scopes, err := s.getUserFunc(username, password)
	return id, scopes, errmgo.Err(err)
}

func (s *Storage) GetClient(cid string) (clover.Client, error) {
	session := s.session.Copy()
	defer session.Close()

	if s.getClientFunc == nil {
		panic("Not implement GetClientFunc")
	}

	client, err := s.getClientFunc(cid)
	return client, errmgo.Err(err)

	// var c clover.DefaultClient
	// if err := session.DB(s.db).C("oauth_client").FindId(cid).One(&c); err != nil {
	// 	return nil, errmgo.Err(err)
	// }
}

func (s *Storage) SetAccessToken(at *clover.AccessToken) error {
	session := s.session.Copy()
	session.SetSafe(nil)
	defer session.Close()

	if err := session.DB(s.db).C("oauth_access_token").Insert(at); err != nil {
		return errmgo.Err(err)
	}

	return nil
}

func (s *Storage) GetAccessToken(at string) (*clover.AccessToken, error) {
	session := s.session.Copy()
	defer session.Close()

	var a *clover.AccessToken
	if err := session.DB(s.db).C("oauth_access_token").FindId(at).One(&a); err != nil {
		return nil, errmgo.Err(err)
	}

	return a, nil
}

func (s *Storage) SetRefreshToken(rt *clover.RefreshToken) error {
	session := s.session.Copy()
	session.SetSafe(nil)
	defer session.Close()

	if err := session.DB(s.db).C("oauth_refresh_token").Insert(rt); err != nil {
		return errmgo.Err(err)
	}

	return nil
}

func (s *Storage) GetRefreshToken(rt string) (*clover.RefreshToken, error) {
	session := s.session.Copy()
	defer session.Close()

	var r *clover.RefreshToken
	if err := session.DB(s.db).C("oauth_refresh_token").FindId(rt).One(&r); err != nil {
		return nil, errmgo.Err(err)
	}

	return r, nil
}

func (s *Storage) RemoveRefreshToken(rt string) error {
	session := s.session.Copy()
	defer session.Close()

	return errmgo.Err(session.DB(s.db).C("oauth_refresh_token").RemoveId(rt))
}

func (s *Storage) SetAuthorizeCode(ac *clover.AuthorizeCode) error {
	session := s.session.Copy()
	session.SetSafe(nil)
	defer session.Close()

	if err := session.DB(s.db).C("oauth_auth_code").Insert(ac); err != nil {
		return errmgo.Err(err)
	}

	return nil
}

func (s *Storage) GetAuthorizeCode(code string) (*clover.AuthorizeCode, error) {
	session := s.session.Copy()
	defer session.Close()

	var ac *clover.AuthorizeCode
	if err := session.DB(s.db).C("oauth_auth_code").FindId(code).One(&ac); err != nil {
		return nil, errmgo.Err(err)
	}

	return ac, nil
}
