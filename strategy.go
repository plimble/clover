package clover

type AuthorizeConfig struct {
	CookieSecret string
	Consent      Consent
	resTypes     map[string]ResponseType
}

func (c *AuthorizeConfig) RegisterResponseType(resType ResponseType) {
	if c.resTypes == nil {
		c.resTypes = make(map[string]ResponseType)
	}

	c.resTypes[resType.Name()] = resType
}

type Strategy struct {
	accessTokenGenerator    TokenGenerator
	refreshTokenGenerator   TokenGenerator
	authorizeTokenGenerator TokenGenerator
	authorizeConfig         *AuthorizeConfig
	grantTypes              map[string]GrantType
	scopeValidator          ScopeValidator
}

func NewStrategy(accessTokenGenerator, refreshTokenGenerator, authorizeTokenGenerator TokenGenerator, scopeValidator ScopeValidator) *Strategy {
	return &Strategy{
		grantTypes: make(map[string]GrantType),
	}
}

func (s *Strategy) AllowAuthorize(config *AuthorizeConfig) {
	s.authorizeConfig = config
}

func (s *Strategy) RegisterGrantType(grant GrantType) {
	s.grantTypes[grant.Name()] = grant
}
