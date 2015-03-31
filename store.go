package clover

//go:generate mockery -name allstore -inpkg
//go:generate msgp

type allstore interface {
	AuthServerStore
	UserStore
	RefreshTokenStore
	AuthCodeStore
	PublicKeyStore
}

type AuthServerStore interface {
	GetClient(id string) (Client, error)
	SetAccessToken(accessToken *AccessToken) error
	GetAccessToken(at string) (*AccessToken, error)
}

type UserStore interface {
	GetUser(username, password string) (string, []string, error)
}

type RefreshTokenStore interface {
	RemoveRefreshToken(rt string) error
	SetRefreshToken(rt *RefreshToken) error
	GetRefreshToken(rt string) (*RefreshToken, error)
}

type AuthCodeStore interface {
	SetAuthorizeCode(ac *AuthorizeCode) error
	GetAuthorizeCode(code string) (*AuthorizeCode, error)
}

type PublicKeyStore interface {
	GetKey(clientID string) (*PublicKey, error)
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
	RefreshToken string   `json:"refresh_token" bson:"_id" msg:"r"`
	ClientID     string   `json:"client_id" bson:"c" msg:"a"`
	UserID       string   `json:"user_id" bson:"u" msg:"u"`
	Expires      int64    `json:"expires" bson:"e" msg:"e"`
	Scope        []string `json:"scope" bson:"s" msg:"s"`
}

type AuthorizeCode struct {
	Code        string   `json:"code" bson:"_id" msg:"co"`
	ClientID    string   `json:"client_id" bson:"c" msg:"c"`
	UserID      string   `json:"user_id" bson:"u" msg:"u"`
	Expires     int64    `json:"expires" bson:"e" msg:"e"`
	Scope       []string `json:"scope" bson:"s" msg:"s"`
	RedirectURI string   `json:"redirect_uri" bson:"r" msg:"r"`
}

type AccessToken struct {
	AccessToken string   `json:"access_token" bson:"_id" msg:"a"`
	ClientID    string   `json:"client_id" bson:"c" msg:"c"`
	UserID      string   `json:"user_id" bson:"u" msg:"u"`
	Expires     int64    `json:"expires" bson:"e" msg:"e"`
	Scope       []string `json:"scope" bson:"s" msg:"s"`
}

type PublicKey struct {
	ClientID   string `json:"client_id" bson:"_id" msg:"c"`
	PublicKey  string `json:"public_key" bson:"pu" msg:"pu"`
	PrivateKey string `json:"private_key" bson:"pr" msg:"pr"`
	Algorithm  string `json:"algorithm" bson:"a" msg:"a"`
}
