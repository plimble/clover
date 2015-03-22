package clover

type AuthConfig struct {
	AccessLifeTime       int
	AuthCodeLifetime     int
	RefreshTokenLifetime int

	DefaultScopes        []string
	AllowCredentialsBody bool
	AllowImplicit        bool
	StateParamRequired   bool

	//store
	AuthServerStore   AuthServerStore
	UserStore         UserStore
	RefreshTokenStore RefreshTokenStore
	AuthCodeStore     AuthCodeStore
	PublicKeyStore    PublicKeyStore

	Grants Grants
}

func NewAuthConfig(authStore AuthServerStore) *AuthConfig {
	return &AuthConfig{
		AccessLifeTime:       3600,
		AuthCodeLifetime:     60,
		RefreshTokenLifetime: 3600,
		AllowCredentialsBody: false,
		AllowImplicit:        false,
		StateParamRequired:   false,
		AuthServerStore:      authStore,
		Grants:               make(Grants),
	}
}

func (c *AuthConfig) UseJWTAccessTokens(store PublicKeyStore) {
	c.PublicKeyStore = store
}

func (c *AuthConfig) AddGrant(grant GrantType) {
	c.Grants[grant.GetGrantType()] = grant
}

func (c *AuthConfig) AddAuthCodeGrant(store AuthCodeStore) {
	c.AuthCodeStore = store
	c.Grants[AUTHORIZATION_CODE] = newAuthCodeGrant(store)
}

func (c *AuthConfig) AddClientGrant() {
	c.Grants[CLIENT_CREDENTIALS] = newClientGrant(c.AuthServerStore)
}

func (c *AuthConfig) AddPasswordGrant(store UserStore) {
	c.UserStore = store
	c.Grants[PASSWORD] = newPasswordGrant(store)
}

func (c *AuthConfig) AddRefreshGrant(store RefreshTokenStore) {
	c.RefreshTokenStore = store
	c.Grants[REFRESH_TOKEN] = newRefreshGrant(store)
}

func (c *AuthConfig) SetDefaultScopes(ids ...string) {
	c.DefaultScopes = ids
}
