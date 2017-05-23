package oauth2

import "time"

type RefreshToken struct {
	RefreshToken         string                 `json:"rt"`
	ClientID             string                 `json:"cid"`
	UserID               string                 `json:"uid"`
	Expired              int64                  `json:"exp"`
	Scopes               []string               `json:"scp"`
	AccessTokenLifespan  int                    `json:"atl"`
	RefreshTokenLifespan int                    `json:"rtl"`
	Extras               map[string]interface{} `json:"ext"`
}

func (r *RefreshToken) Valid() bool {
	return r != nil && r.RefreshToken != "" && time.Now().UTC().Unix() > r.Expired
}
