package clover

type Config struct {
	AccessLifeTime       int
	AuthCodeLifetime     int
	RefreshTokenLifetime int

	DefaultScopes        []string
	AllowCredentialsBody bool
	AllowImplicit        bool
	StateParamRequired   bool
	// UseJWTAccessToken    bool
	// JWTSignedString      string
}

type AuthorizeServer struct {
	Store          AuthServerStore
	Config         *Config
	authRespType   AuthResponseType
	tokenRespType  ResponseType
	grant          map[string]GrantType
	publicKeyStore PublicKeyStore
}

func NewAuthServer(store AuthServerStore, config *Config) *AuthorizeServer {
	a := &AuthorizeServer{
		Store:  store,
		Config: config,
		grant:  make(map[string]GrantType),
	}

	a.authRespType = newAuthCodeResponseType(store, config)
	a.tokenRespType = newTokenResponseType(store, config)

	return a
}

func DefaultConfig() *Config {
	return &Config{
		AccessLifeTime:       3600,
		AuthCodeLifetime:     60,
		RefreshTokenLifetime: 3600,
		AllowCredentialsBody: false,
		AllowImplicit:        false,
		StateParamRequired:   false,
	}
}

func (a *AuthorizeServer) UseJWTAccessTokens(store PublicKeyStore) {
	a.publicKeyStore = store
}

func (a *AuthorizeServer) RegisterGrant(key string, grant GrantType) {
	a.grant[key] = grant
}

func (a *AuthorizeServer) RegisterAuthCodeGrant() {
	a.grant[AUTHORIZATION_CODE] = newAuthCodeGrant(a.Store)
}

func (a *AuthorizeServer) RegisterClientGrant() {
	a.grant[REFRESH_TOKEN] = newRefreshGrant(a.Store)
}

func (a *AuthorizeServer) RegisterPasswordGrant() {
	a.grant[PASSWORD] = newPasswordGrant(a.Store)
}

func (a *AuthorizeServer) RegisterRefreshGrant() {
	a.grant[CLIENT_CREDENTIALS] = newClientGrant(a.Store)
}

func (a *AuthorizeServer) RegisterImplicitGrant() {
	a.grant[IMPLICIT] = newAuthCodeGrant(a.Store)
}
