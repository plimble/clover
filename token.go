package clover

import (
	"strings"
)

type TokenRequest struct {
	Username     string
	Password     string
	Code         string
	RedirectURI  string
	RefreshToken string
	ClientID     string
	ClientSecret string
	GrantType    string
	Scope        string
	Assertion    string
}

type TokenData struct {
	GrantType GrantType
	GrantData *GrantData
	Scope     []string
}

type tokenController struct {
	config        *Config
	store         AuthServerStore
	grant         Grants
	tokenRespType AccessTokenResponseType
}

func newTokenController(config *Config, store AuthServerStore, grant Grants, tokenRespType AccessTokenResponseType) *tokenController {
	return &tokenController{
		config:        config,
		store:         store,
		grant:         grant,
		tokenRespType: tokenRespType,
	}
}

func (t *tokenController) handleAccessToken(tr *TokenRequest) *response {
	td := &TokenData{}
	if resp := t.validateAccessToken(tr, td); resp != nil {
		return resp
	}

	return td.GrantType.CreateAccessToken(td, t.tokenRespType)
}

func (t *tokenController) validateAccessToken(tr *TokenRequest, td *TokenData) *response {
	var resp *response
	if td.GrantType, resp = t.validateAccessTokenGrantType(tr); resp != nil {
		return resp
	}

	if td.GrantData, resp = td.GrantType.Validate(tr); resp != nil {
		return resp
	}

	if td.GrantType.GetGrantType() != CLIENT_CREDENTIALS {
		if resp = t.validateAccessTokenClient(tr, td.GrantData); resp != nil {
			return resp
		}
	}

	if !checkGrantType(td.GrantData.GrantType, td.GrantType.GetGrantType()) {
		return errUnAuthorizedGrant
	}

	if td.Scope, resp = t.validateAccessTokenScope(tr, td.GrantData); resp != nil {
		return resp
	}

	return nil
}

func (t *tokenController) validateAccessTokenGrantType(tr *TokenRequest) (GrantType, *response) {
	if tr.GrantType == "" {
		return nil, errGrantTypeRequired
	}

	if _, ok := t.grant[tr.GrantType]; !ok {
		return nil, errUnSupportedGrantType
	}

	return t.grant[tr.GrantType], nil
}

func (t *tokenController) validateAccessTokenClient(tr *TokenRequest, grantData *GrantData) *response {
	var err error
	client, err := t.store.GetClient(tr.ClientID)
	if err != nil {
		return errInternal(err.Error())
	}

	if grantData.ClientID != "" && grantData.ClientID != client.GetClientID() {
		return errInvalidClientCredentail
	}

	copyClientToGrant(grantData, client)

	return nil
}

func (t *tokenController) validateAccessTokenScope(tr *TokenRequest, grantData *GrantData) ([]string, *response) {
	scopes := strings.Fields(tr.Scope)
	if len(scopes) > 0 {
		if grantData.Scope != nil && len(grantData.Scope) > 0 {
			if !checkScope(grantData.Scope, scopes...) {
				return nil, errInvalidScopeRequest
			}
		} else {
			return nil, errNoScope
		}
	} else if grantData.Scope != nil && len(grantData.Scope) > 0 {
		scopes = grantData.Scope
	} else {
		if len(t.config.DefaultScopes) == 0 {
			return nil, errNoScope
		}

		scopes = t.config.DefaultScopes
	}

	return scopes, nil
}
