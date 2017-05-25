package oauth2

import "time"

type AccessToken struct {
	AccessToken string                 `json:"at"`
	ClientID    string                 `json:"cid,omitempty"`
	UserID      string                 `json:"uid,omitempty"`
	Expired     int64                  `json:"exp"`
	ExpiresIn   int                    `json:"ein"`
	Scopes      []string               `json:"scp"`
	Extras      map[string]interface{} `json:"ext"`
}

func (a *AccessToken) Valid() bool {
	return a != nil && a.AccessToken != "" && time.Now().UTC().Unix() > a.Expired
}

func (a *AccessToken) HasScope(scopes ...string) bool {
	for _, scope := range scopes {
		if ok := HierarchicScope(scope, a.Scopes); !ok {
			return false
		}
	}

	return true
}

type JWTAccessToken struct {
	Audience  string
	ExpiresAt int64
	ID        string
	IssuedAt  int64
	Issuer    string
	Subject   string
	Extras    map[string]interface{}
	Scopes    []string
}

func (a *JWTAccessToken) Valid() bool {
	return a != nil && time.Now().UTC().Unix() > int64(a.ExpiresAt)
}

func (a *JWTAccessToken) HasScope(scopes ...string) bool {
	for _, scope := range scopes {
		if ok := HierarchicScope(scope, a.Scopes); !ok {
			return false
		}
	}

	return true
}
