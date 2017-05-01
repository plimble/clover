package clover

import "strings"
import "fmt"

//go:generate mockery -name OAuth2 -case underscore
type OAuth2 interface {
	RegisterGrantType(grant GrantType)
}

type oauth2 struct {
	config       *Config
	accessStore  AccessTokenStore
	refreshStore RefreshTokenStore
	client       ClientManager
	grantTypes   map[string]GrantType
}

func New(config *Config, client ClientManager, accessStore AccessTokenStore, refreshStore RefreshTokenStore) OAuth2 {
	return &oauth2{
		config:       config,
		accessStore:  accessStore,
		refreshStore: refreshStore,
		client:       client,
		grantTypes:   make(map[string]GrantType),
	}
}

func (o *oauth2) RegisterGrantType(grant GrantType) {
	o.grantTypes[grant.Name()] = grant
}

type TokenRes struct {
}

func (o *oauth2) Token(ctx Context) (*TokenRes, error) {
	if ctx.GrantType == "" {
		return nil, ErrGrantTypeNotFound()
	}

	if ctx.ClientID == "" || ctx.ClientSecret == "" {
		return nil, ErrClientCredentialRequired()
	}

	grant, ok := o.grantTypes[ctx.GrantType]
	if !ok {
		return nil, ErrGrantTypeNotSupport(grant.Name())
	}

	_, err := o.client.GetClient(ctx.ClientID, ctx.ClientSecret)
	if err != nil {
		return nil, ErrInvalidClient(ctx.ClientID, err)
	}

	var grantData *GrantData
	if grantData, err = grant.Validate(ctx); err != nil {
		return nil, err
	}

	var scopes []string
	ctxScopes := strings.Fields(ctx.Scope)

	if len(grantData.Scopes) == 0 {
		return nil, ErrUnSupportedScope()
	}

	switch {
	case ctx.Scope != "":
		if len(grantData.Scopes) > 0 {
			if ok = CheckScope(grantData.Scopes, ctxScopes); !ok {
				return nil, ErrInvalidScope()
			}
			scopes = ctxScopes
		}
	default:
		scopes = grantData.Scopes
	}

	fmt.Println(scopes)

	// create access token and return response

	return nil, nil
}

type AuthRes struct {
}

func (o *oauth2) Auth(ctx Context) (*AuthRes, error) {
	return nil, nil
}

type RevokeRes struct {
}

func (o *oauth2) Revoke(ctx Context) (*RevokeRes, error) {
	return nil, nil
}
