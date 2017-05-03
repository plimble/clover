package clover

type AccessToken struct {
	AccessToken string   `json:"access_token"`
	ClientID    string   `json:"client_id"`
	UserID      string   `json:"user_id"`
	Expired     int64    `json:"expired"`
	Scopes      []string `json:"scopes"`
}

type RefreshToken struct {
	RefreshToken string   `json:"refresh_token"`
	ClientID     string   `json:"client_id"`
	UserID       string   `json:"user_id"`
	Expired      int64    `json:"expired"`
	Scopes       []string `json:"scopes"`
}

type AuthorizeCode struct {
	Code        string   `json:"code"`
	ClientID    string   `json:"client_id"`
	UserID      string   `json:"user_id"`
	Expired     int64    `json:"expired"`
	Scopes      []string `json:"scopes"`
	RedirectURI string   `json:"redirect_uri"`
}

type Session struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Expired   int64                  `json:"expired"`
	ExtraData map[string]interface{} `json:"extra_data"`
}

type User struct {
	ID        string                 `json:"id"`
	Username  string                 `json:"username"`
	Password  string                 `json:"password"`
	ExtraData map[string]interface{} `json:"extra_data"`
}

//go:generate mockery -name ClientStore -case underscore
type ClientStore interface {
	DeleteClient(id string) error
	SaveClient(client *Client) error
	GetClientWithSecret(id, secret string) (*Client, error)
	GetClient(id string) (*Client, error)
}

//go:generate mockery -name SessionStore -case underscore
type SessionStore interface {
	DeleteSession(id string) error
	SaveSession(session *Session) error
	GetSession(id string) (*Session, error)
}

//go:generate mockery -name UserStore -case underscore
type UserStore interface {
	GetUser(username, password string) (*User, error)
}

//go:generate mockery -name TokenStore -case underscore
type TokenStore interface {
	DeleteAccessToken(token string) error
	SaveAccessToken(accessToken *AccessToken) error
	GetAccessToken(at string) (*AccessToken, error)

	DeleteRefreshToken(token string) error
	SaveRefreshToken(token *RefreshToken) error
	GetRefreshToken(token string) (*RefreshToken, error)

	DeleteAuthCode(code string) error
	SaveAuthCode(code *AuthorizeCode) error
	GetAuthCode(code string) (*AuthorizeCode, error)
}
