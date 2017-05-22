package oauth2

import "time"

type AccessToken struct {
	AccessToken string                 `json:"access_token"`
	ClientID    string                 `json:"client_id,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	Expired     int64                  `json:"expired"`
	ExpiresIn   int                    `json:"expires_in"`
	Scopes      []string               `json:"scopes"`
	Extras      map[string]interface{} `json:"extras"`
}

func (a *AccessToken) Valid() bool {
	return a != nil && a.AccessToken != "" && time.Now().UTC().Unix() > a.Expired
}
