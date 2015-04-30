package clover

//go:generate mockery -name allstore -inpkg
//go:generate msgp

type allstore interface {
	ClientStore
	UserStore
	AccessTokenStore
	RefreshTokenStore
	AuthCodeStore
	PublicKeyStore
	ScopeStore
}

type ClientStore interface {
	GetClient(id string) (Client, error)
}

type ScopeStore interface {
	ExistScopes(scopes ...string) (bool, error)
	GetDefaultScope(clientID string) ([]string, error)
}

type UserStore interface {
	GetUser(username, password string) (User, error)
}

type AccessTokenStore interface {
	SetAccessToken(accessToken *AccessToken) error
	GetAccessToken(at string) (*AccessToken, error)
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

type User interface {
	GetID() string
	GetUsername() string
	GetPassword() string
	GetData() map[string]interface{}
	GetScope() []string
}

type Client interface {
	GetClientID() string
	GetClientSecret() string
	GetGrantType() []string
	GetUserID() string
	GetScope() []string
	GetRedirectURI() string
	GetData() map[string]interface{}
}

type RefreshToken struct {
	RefreshToken string                 `json:"refresh_token" bson:"_id" msg:"r"`
	ClientID     string                 `json:"client_id" bson:"c" msg:"a"`
	UserID       string                 `json:"user_id" bson:"u" msg:"u"`
	Expires      int64                  `json:"expires" bson:"e" msg:"e"`
	Scope        []string               `json:"scope" bson:"s" msg:"s"`
	Data         map[string]interface{} `json:"data" bson:"d" msg:"d"`
}

type AuthorizeCode struct {
	Code        string                 `json:"code" bson:"_id" msg:"co"`
	ClientID    string                 `json:"client_id" bson:"c" msg:"c"`
	UserID      string                 `json:"user_id" bson:"u" msg:"u"`
	Expires     int64                  `json:"expires" bson:"e" msg:"e"`
	Scope       []string               `json:"scope" bson:"s" msg:"s"`
	RedirectURI string                 `json:"redirect_uri" bson:"r" msg:"r"`
	Data        map[string]interface{} `json:"data" bson:"d" msg:"d"`
}

type AccessToken struct {
	AccessToken string                 `json:"access_token" bson:"_id" msg:"a"`
	ClientID    string                 `json:"client_id" bson:"c" msg:"c"`
	UserID      string                 `json:"user_id" bson:"u" msg:"u"`
	Expires     int64                  `json:"expires" bson:"e" msg:"e"`
	Scope       []string               `json:"scope" bson:"s" msg:"s"`
	Data        map[string]interface{} `json:"data" bson:"d" msg:"d"`
}

type PublicKey struct {
	ClientID   string `json:"client_id" bson:"_id" msg:"c"`
	PublicKey  string `json:"public_key" bson:"pu" msg:"pu"`
	PrivateKey string `json:"private_key" bson:"pr" msg:"pr"`
	Algorithm  string `json:"algorithm" bson:"a" msg:"a"`
}
