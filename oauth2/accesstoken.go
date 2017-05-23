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
