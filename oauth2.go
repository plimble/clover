package clover

//go:generate mockery -name OAuth2 -case underscore
type OAuth2 interface {
	RegisterGrantType(grant GrantType)
}

type oauth2 struct {
	config        *Config
	clientManager ClientManager
	sessionStore  SessionStore
	tokenManager  TokenManager
	grantTypes    map[string]GrantType
}

func New(config *Config, clientManager ClientManager, tokenManager TokenManager, sessionStore SessionStore) OAuth2 {
	return &oauth2{
		config:        config,
		clientManager: clientManager,
		sessionStore:  sessionStore,
		tokenManager:  tokenManager,
		grantTypes:    make(map[string]GrantType),
	}
}

func (o *oauth2) RegisterGrantType(grant GrantType) {
	o.grantTypes[grant.Name()] = grant
}
