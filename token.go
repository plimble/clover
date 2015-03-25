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
	config        *AuthConfig
	tokenRespType TokenRespType
}

func newTokenController(config *AuthConfig, tokenRespType TokenRespType) *tokenController {
	return &tokenController{
		config:        config,
		tokenRespType: tokenRespType,
	}
}

func (t *tokenController) handleAccessToken(tr *TokenRequest) *Response {
	td := &TokenData{}
	if resp := t.validateToken(tr, td); resp != nil {
		return resp
	}

	return td.GrantType.CreateAccessToken(td, t.tokenRespType)
}

func (t *tokenController) validateToken(tr *TokenRequest, td *TokenData) *Response {
	var resp *Response
	if td.GrantType, resp = t.validateGrantType(tr); resp != nil {
		return resp
	}

	if td.GrantData, resp = td.GrantType.Validate(tr); resp != nil {
		return resp
	}

	if td.GrantType.GetGrantType() != CLIENT_CREDENTIALS {
		if resp = t.validateClient(tr, td.GrantData); resp != nil {
			return resp
		}
	}

	if !checkGrantType(td.GrantData.GrantType, td.GrantType.GetGrantType()) {
		return errUnAuthorizedGrant
	}

	if td.Scope, resp = t.validateScope(tr, td.GrantData); resp != nil {
		return resp
	}

	return nil
}

func (t *tokenController) validateGrantType(tr *TokenRequest) (GrantType, *Response) {
	if tr.GrantType == "" {
		return nil, errGrantTypeRequired
	}

	if _, ok := t.config.Grants[tr.GrantType]; !ok {
		return nil, errUnSupportedGrantType
	}

	return t.config.Grants[tr.GrantType], nil
}

func (t *tokenController) validateClient(tr *TokenRequest, grantData *GrantData) *Response {
	var err error
	client, err := t.config.AuthServerStore.GetClient(tr.ClientID)
	if err != nil {
		return errInvalidClientID
	}

	if grantData.ClientID != "" && grantData.ClientID != client.GetClientID() {
		return errInvalidClientCredentail
	}

	copyClientToGrant(grantData, client)

	return nil
}

func (t *tokenController) validateScope(tr *TokenRequest, grantData *GrantData) ([]string, *Response) {
	scopes := strings.Fields(tr.Scope)
	if len(scopes) > 0 {
		if len(grantData.Scope) > 0 {
			if !checkScope(grantData.Scope, scopes...) {
				return nil, errInvalidScopeRequest
			}
		} else {
			return nil, errUnSupportedScope
		}
	} else if len(grantData.Scope) > 0 {
		scopes = grantData.Scope
	} else {
		scopes = t.config.DefaultScopes
	}

	return scopes, nil
}
