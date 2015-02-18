package clover

type Config struct {
	//Store
	Store                Store
	AccessLifeTime       int
	AuthCodeLifetime     int
	RefreshTokenLifetime int
	Grants               []string

	DefaultScopes        []string
	AllowCredentialsBody bool
	AllowImplicit        bool
	StateParamRequired   bool
	// UseJWTAccessToken    bool
	// JWTSignedString      string
}

func DefaultConfig() *Config {
	return &Config{
		AccessLifeTime:       3600,
		AuthCodeLifetime:     30,
		RefreshTokenLifetime: 1209600,
		AllowCredentialsBody: true,
	}
}

//main package endpoint
type Clover struct {
	Config    *Config
	tokenType TokenType
	RespType  map[string]ResponseType
	Grant     map[string]GrantType
}

func New(config *Config) *Clover {
	c := &Clover{}
	c.Config = config
	c.tokenType = newBearerTokenType()
	c.RespType = map[string]ResponseType{
		"token": newTokenResponseType(),
		"code":  newCodeResponseType(),
	}
	c.Grant = map[string]GrantType{
		AUTHORIZATION_CODE: newAuthCodeGrant(),
		REFRESH_TOKEN:      newRefreshGrant(),
		PASSWORD:           newPasswordGrant(),
		CLIENT_CREDENTIALS: newClientGrant(),
	}

	return c
}
