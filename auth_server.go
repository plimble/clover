package clover

type Config struct {
	//Store
	Store                Store
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
	Config   *Config
	RespType map[string]ResponseType
	Grant    map[string]GrantType
}

func NewAuthorizeServer(config *Config) *AuthorizeServer {
	if config.Store == nil {
		panic("store should not be nil")
	}

	return &AuthorizeServer{
		Config: config,
		RespType: map[string]ResponseType{
			RESP_TYPE_TOKEN: newTokenResponseType(),
		},
		Grant: make(map[string]GrantType),
	}
}
