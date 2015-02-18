package mongo

import (
	"github.com/plimble/clover"
	"github.com/plimble/errs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct {
	session     *mgo.Session
	db          string
	GetUserFunc GetUserFunc
}

type GetUserFunc func(session *mgo.Session, username, password string) (string, error)

func New(session *mgo.Session, db string) *Storage {
	return &Storage{session, db, nil}
}

func (s *Storage) TruncateAll() {
	session := s.session.Copy()
	defer session.Close()
	// session.DB(s.db).C("activity").RemoveAll(nil)
}

func (s *Storage) GetUser(username, password string) (string, error) {
	session := s.session.Copy()
	defer session.Close()

	if s.GetUser == nil {
		panic("Not implement GetUserFunc")
	}

	return s.GetUserFunc(session, username, password)
}

func (s *Storage) GetClient(cid string) (*clover.Client, error) {
	session := s.session.Copy()
	defer session.Close()

	var c clover.Client
	if err := session.DB(s.db).C("oauth_client").FindId(cid).One(&c); err != nil {
		return nil, errs.Mgo(err)
	}

	return &c, nil
}

func (s *Storage) SetAccessToken(at *clover.AccessToken) error {
	session := s.session.Copy()
	session.SetSafe(nil)
	defer session.Close()

	if err := session.DB(s.db).C("oauth_access_token").Insert(at); err != nil {
		return errs.Mgo(err)
	}

	return nil
}

func (s *Storage) GetAccessToken(at string) (*clover.AccessToken, error) {
	session := s.session.Copy()
	defer session.Close()

	var a *clover.AccessToken
	if err := session.DB(s.db).C("oauth_access_token").FindId(at).One(&a); err != nil {
		return nil, errs.Mgo(err)
	}

	return a, nil
}

func (s *Storage) SetRefreshToken(rt *clover.RefreshToken) error {
	session := s.session.Copy()
	session.SetSafe(nil)
	defer session.Close()

	if err := session.DB(s.db).C("oauth_refresh_token").Insert(rt); err != nil {
		return errs.Mgo(err)
	}

	return nil
}

func (s *Storage) GetRefreshToken(rt string) (*clover.RefreshToken, error) {
	session := s.session.Copy()
	defer session.Close()

	var r *clover.RefreshToken
	if err := session.DB(s.db).C("oauth_refresh_token").FindId(rt).One(&r); err != nil {
		return nil, errs.Mgo(err)
	}

	return r, nil
}

func (s *Storage) RemoveRefreshToken(rt string) error {
	session := s.session.Copy()
	defer session.Close()

	return errs.Mgo(session.DB(s.db).C("oauth_refresh_token").RemoveId(rt))
}

func (s *Storage) SetAuthorizeCode(ac *clover.AuthorizeCode) error {
	session := s.session.Copy()
	session.SetSafe(nil)
	defer session.Close()

	if err := session.DB(s.db).C("oauth_auth_code").Insert(ac); err != nil {
		return errs.Mgo(err)
	}

	return nil
}

func (s *Storage) GetAuthorizeCode(code string) (*clover.AuthorizeCode, error) {
	session := s.session.Copy()
	defer session.Close()

	var ac *clover.AuthorizeCode
	if err := session.DB(s.db).C("oauth_auth_code").FindId(code).One(&ac); err != nil {
		return nil, errs.Mgo(err)
	}

	return ac, nil
}

func (s *Storage) SetScope(id, desc string) error {
	session := s.session.Copy()
	session.SetSafe(nil)
	defer session.Close()

	if err := session.DB(s.db).C("oauth_scope").Insert(clover.Scope{id, desc}); err != nil {
		return errs.Mgo(err)
	}

	return nil
}

func (s *Storage) GetScopes(ids []string) ([]*clover.Scope, error) {
	session := s.session.Copy()
	defer session.Close()

	var scope []*clover.Scope
	if err := session.DB(s.db).C("oauth_scope").Find(bson.M{"_id": bson.M{"$in": ids}}).All(&scope); err != nil {
		return nil, errs.Mgo(err)
	}

	return scope, nil
}

func (s *Storage) GetAllScopeID() ([]string, error) {
	session := s.session.Copy()
	defer session.Close()

	var scopeID []string
	var scope *clover.Scope

	iter := session.DB(s.db).C("oauth_scope").Find(nil).Select(bson.M{"_id": 1}).Iter()
	for iter.Next(&scope) {
		scopeID = append(scopeID, scope.ID)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return scopeID, nil
}
