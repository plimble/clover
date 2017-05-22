package oauth2

import "time"

type RefreshToken struct {
	RefreshToken         string                 `json:"refresh_token"`
	ClientID             string                 `json:"client_id"`
	UserID               string                 `json:"user_id"`
	Expired              int64                  `json:"expired"`
	Scopes               []string               `json:"scope"`
	AccessTokenLifespan  int                    `json:"access_token_lifespan"`
	RefreshTokenLifespan int                    `json:"refresh_token_lifespan"`
	Extras               map[string]interface{} `json:"extras"`
}

func (r *RefreshToken) Valid() bool {
	return r != nil && r.RefreshToken != "" && time.Now().UTC().Unix() > r.Expired
}
