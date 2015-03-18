package clover

type Store interface {
	GetClient(id string) (Client, error)

	SetAccessToken(accessToken *AccessToken) error
	GetAccessToken(at string) (*AccessToken, error)

	SetRefreshToken(rt *RefreshToken) error
	GetRefreshToken(rt string) (*RefreshToken, error)
	RemoveRefreshToken(rt string) error

	SetAuthorizeCode(ac *AuthorizeCode) error
	GetAuthorizeCode(code string) (*AuthorizeCode, error)

	GetUser(username, password string) (string, []string, error)
}

type Client interface {
	GetClientID() string
	GetClientSecret() string
	GetGrantType() []string
	GetUserID() string
	GetScope() []string
	GetRedirectURI() string
}

type DefaultClient struct {
	ClientID     string
	ClientSecret string
	GrantType    []string
	UserID       string
	Scope        []string
	RedirectURI  string
}

func (c *DefaultClient) GetClientID() string {
	return c.ClientID
}

func (c *DefaultClient) GetClientSecret() string {
	return c.ClientSecret
}

func (c *DefaultClient) GetGrantType() []string {
	return c.GrantType
}

func (c *DefaultClient) GetUserID() string {
	return c.UserID
}

func (c *DefaultClient) GetScope() []string {
	return c.Scope
}

func (c *DefaultClient) GetRedirectURI() string {
	return c.RedirectURI
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
