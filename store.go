package clover

type Store interface {
	GetClient(id string) (*Client, error)

	SetAccessToken(accessToken *AccessToken) error
	GetAccessToken(at string) (*AccessToken, error)

	SetRefreshToken(rt *RefreshToken) error
	GetRefreshToken(rt string) (*RefreshToken, error)
	RemoveRefreshToken(rt string) error

	SetAuthorizeCode(ac *AuthorizeCode) error
	GetAuthorizeCode(code string) (*AuthorizeCode, error)

	SetScope(id, desc string) error
	GetScopes(ids []string) ([]*Scope, error)
	GetAllScopeID() ([]string, error)

	GetUser(username, password string) (string, error)
}

type Client struct {
	ClientID     string   `json:"client_id" bson:"_id"`
	ClientSecret string   `json:"client_secret" bson:"client_secret"`
	GrantType    []string `json:"grant_type" bson:"grant_type"`
	UserID       string   `json:"user_id" bson:"user_id"`
	Scope        []string `json:"scope" bson:"scope"`
	RedirectURI  string   `json:"redirect_uri" bson:"redirect_uri"`
}

type RefreshToken struct {
	RefreshToken string   `json:"refresh_token" bson:"_id"`
	ClientID     string   `json:"client_id" bson:"client_id"`
	UserID       string   `json:"user_id" bson:"user_id"`
	Expires      int64    `json:"expires" bson:"expires"`
	Scope        []string `json:"scope" bson:"scope"`
}

type AuthorizeCode struct {
	Code        string   `json:"code" bson:"_id"`
	ClientID    string   `json:"client_id" bson:"client_id"`
	UserID      string   `json:"user_id" bson:"user_id"`
	Expires     int64    `json:"expires" bson:"expires"`
	Scope       []string `json:"scope" bson:"scope"`
	RedirectURI string   `json:"redirect_uri" bson:"redirect_uri"`
}

type AccessToken struct {
	AccessToken string   `json:"access_token" bson:"_id"`
	ClientID    string   `json:"client_id" bson:"client_id"`
	UserID      string   `json:"user_id" bson:"user_id"`
	Expires     int64    `json:"expires" bson:"expires"`
	Scope       []string `json:"scope" bson:"scope"`
}

type Scope struct {
	ID   string `json:"id" bson:"_id"`
	Desc string `json:"desc" bson:"desc"`
}
