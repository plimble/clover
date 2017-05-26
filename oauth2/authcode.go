package oauth2

import "time"

type AuthorizeCode struct {
	Code         string                 `json:"c"`
	ClientID     string                 `json:"cid"`
	UserID       string                 `json:"uid"`
	Expired      int64                  `json:"exp"`
	Scopes       []string               `json:"scp"`
	RedirectURI  string                 `json:"rdr"`
	ResponseType string                 `json:"rpt"`
	Extras       map[string]interface{} `json:"ext"`
}

func (a *AuthorizeCode) Valid() bool {
	return a != nil && a.Code != "" && time.Now().UTC().Unix() < a.Expired
}
