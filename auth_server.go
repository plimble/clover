package clover

type AuthorizeConfig struct {
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
	Config   *AuthorizeConfig
	RespType map[string]ResponseType
	Grant    map[string]GrantType
}

func NewAuthorizeServer(config *AuthorizeConfig) *AuthorizeServer {
	a := &AuthorizeServer{}
	a.Config = config
	a.RespType = map[string]ResponseType{
		RESP_TYPE_TOKEN: newTokenResponseType(),
		RESP_TYPE_CODE:  newCodeResponseType(),
	}
	a.Grant = make(map[string]GrantType)

	return a
}
