package clover

import (
	"strings"
)

type TokenData struct {
	grantType GrantType
	grantData *GrantData
	client    Client
	Scope     []string
	ClientID  string
	UserID    string
	Data      map[string]interface{}
}

type tokenCtrl struct {
	store  ClientStore
	config *AuthServerConfig
}

func newTokenCtrl(store ClientStore, config *AuthServerConfig) *tokenCtrl {
	return &tokenCtrl{store, config}
}

func (t *tokenCtrl) token(tr *TokenRequest, respType AccessTokenRespType, grants map[string]GrantType) *Response {
	if respType == nil {
		return errInvalidRespType
	}

	td, resp := t.validate(tr, grants)
	if resp != nil {
		return resp
	}

	if resp = td.grantType.BeforeCreateAccessToken(tr, td); resp != nil {
		return resp
	}

	return respType.Response(td, td.grantType.IncludeRefreshToken())
}

func (t *tokenCtrl) validate(tr *TokenRequest, grants map[string]GrantType) (*TokenData, *Response) {
	var resp *Response
	td := &TokenData{}

	if td.grantType, resp = t.validateGrantType(tr, grants); resp != nil {
		return nil, resp
	}

	if td.grantData, resp = td.grantType.Validate(tr); resp != nil {
		return nil, resp
	}

	if td.client, resp = t.validateClient(tr, td.grantData); resp != nil {
		return nil, resp
	}

	if !checkGrantType(td.client.GetGrantType(), td.grantType.Name()) {
		return nil, errUnAuthorizedGrant
	}

	if td.Scope, resp = t.validateScope(tr, td.grantData, td.client); resp != nil {
		return nil, resp
	}

	td.ClientID = td.client.GetClientID()
	if td.grantData.UserID != "" {
		td.UserID = td.grantData.UserID
	}
	if len(td.grantData.Data) > 0 {
		td.Data = td.grantData.Data
	}

	return td, nil
}

func (t *tokenCtrl) validateGrantType(tr *TokenRequest, grants map[string]GrantType) (GrantType, *Response) {
	if tr.GrantType == "" {
		return nil, errGrantTypeRequired
	}

	grant, ok := grants[tr.GrantType]
	if !ok {
		return nil, errUnSupportedGrantType
	}

	return grant, nil
}

func (t *tokenCtrl) validateClient(tr *TokenRequest, grantData *GrantData) (Client, *Response) {
	var err error
	client, err := t.store.GetClient(tr.ClientID)
	if err != nil {
		return nil, errInvalidClientID
	}

	if grantData.ClientID != "" && grantData.ClientID != client.GetClientID() {
		return nil, errInvalidClientCredentail
	}

	return client, nil
}

func (t *tokenCtrl) validateScope(tr *TokenRequest, grantData *GrantData, client Client) ([]string, *Response) {
	scopes := strings.Fields(tr.Scope)
	if len(scopes) > 0 {
		if len(grantData.Scope) > 0 {
			if !checkScope(grantData.Scope, scopes...) {
				return nil, errInvalidScopeRequest
			}
		} else if len(client.GetScope()) > 0 {
			if !checkScope(client.GetScope(), scopes...) {
				return nil, errInvalidScopeRequest
			}
		} else {
			return nil, errUnSupportedScope
		}
	} else if len(grantData.Scope) > 0 {
		scopes = grantData.Scope
	} else if len(client.GetScope()) > 0 {
		scopes = client.GetScope()
	} else {
		if len(t.config.DefaultScopes) == 0 {
			return nil, errUnSupportedScope
		}

		scopes = t.config.DefaultScopes
	}

	return scopes, nil
}
