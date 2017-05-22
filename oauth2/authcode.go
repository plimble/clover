package oauth2

import "time"

type AuthorizeCode struct {
	Code         string                 `json:"code"`
	ClientID     string                 `json:"client_id"`
	UserID       string                 `json:"user_id"`
	Expired      int64                  `json:"expired"`
	Scopes       []string               `json:"scopes"`
	RedirectURI  string                 `json:"redirect_uri"`
	ResponseType string                 `json:"response_type"`
	Extras       map[string]interface{} `json:"extras"`
}

func (a *AuthorizeCode) Valid() bool {
	return a != nil && a.Code != "" && time.Now().UTC().Unix() > a.Expired
}
