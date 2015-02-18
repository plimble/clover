package clover

import (
	"net/http"
	"strings"
)

type tokenRequest struct {
	username     string
	password     string
	code         string
	redirectURI  string
	refreshToken string
	clientID     string
	clientSecret string
	grantType    string
	scope        string
	assertion    string
}

type TokenData struct {
	GrantType GrantType
	GrantData *GrantData
	Client    *Client
	Scope     []string
}

func parseTokenRequest(r *http.Request, c *Clover) (*tokenRequest, *Response) {
	clientID, clientSecret, resp := getCredentialsFromHttp(r, c.Config)
	if resp != nil {
		return nil, resp
	}

	tr := &tokenRequest{
		clientID:     clientID,
		clientSecret: clientSecret,
		scope:        r.FormValue("scope"),
		grantType:    r.FormValue("grant_type"),
		username:     r.FormValue("username"),
		password:     r.FormValue("password"),
		code:         r.FormValue("code"),
		redirectURI:  r.FormValue("redirect_uri"),
		refreshToken: r.FormValue("refresh_token"),
		assertion:    r.FormValue("assertion"),
	}

	return tr, nil
}

func (c *Clover) Token(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		errMustBePostMetthod.Write(w)
		return
	}

	tr, resp := parseTokenRequest(r, c)
	if resp != nil {
		resp.Write(w)
		return
	}

	resp = c.handleAccessToken(tr)
	resp.Write(w)
}

func (c *Clover) handleAccessToken(tr *tokenRequest) *Response {
	td := &TokenData{}
	if resp := c.validateAccessToken(tr, td); resp != nil {
		return resp
	}

	return td.GrantType.CreateAccessToken(td, c, c.RespType[RESP_TYPE_TOKEN])
}

func (c *Clover) validateAccessToken(tr *tokenRequest, td *TokenData) *Response {
	var resp *Response
	if resp := c.validateAccessTokenGrantType(tr, td); resp != nil {
		return resp
	}

	td.GrantData, resp = td.GrantType.Validate(tr, c)
	if resp != nil {
		return resp
	}

	if td.GrantType.GetGrantType() != CLIENT_CREDENTIALS {
		if resp := c.validateAccessTokenClient(tr, td); resp != nil {
			return resp
		}
	} else {
		td.Client = &Client{
			ClientID:  td.GrantData.ClientID,
			UserID:    td.GrantData.UserID,
			Scope:     td.GrantData.Scope,
			GrantType: td.GrantData.GrantType,
		}
	}

	if !checkGrantType(td.Client.GrantType, td.GrantType.GetGrantType()) {
		return errUnAuthorizedGrant
	}

	if resp := c.validateAccessTokenScope(tr, td); resp != nil {
		return resp
	}

	return nil
}

func (c *Clover) validateAccessTokenGrantType(tr *tokenRequest, td *TokenData) *Response {
	if tr.grantType == "" {
		return errGrantTypeRequired
	}

	for i := 0; i < len(c.Config.Grants); i++ {
		if c.Config.Grants[i] == tr.grantType {
			td.GrantType = c.Grant[tr.grantType]
			return nil
		}
	}

	return errUnSupportedGrantType
}

func (c *Clover) validateAccessTokenClient(tr *tokenRequest, td *TokenData) *Response {
	var err error
	td.Client, err = c.Config.Store.GetClient(td.GrantData.ClientID)
	if err != nil {
		return errInternal(err.Error())
	}

	if td.Client.ClientID != "" && td.GrantData.ClientID != td.Client.ClientID {
		return errInvalidClientCredentail
	}

	return nil
}

func (c *Clover) validateAccessTokenScope(tr *tokenRequest, td *TokenData) *Response {
	td.Scope = strings.Fields(tr.scope)
	if len(td.Scope) > 0 {
		if len(td.GrantData.Scope) > 0 {
			if !checkScope(td.Scope, td.GrantData.Scope) {
				return errInvalidScopeRequest
			}
		} else {
			if !checkScope(td.Scope, td.Client.Scope) {
				return errInvalidScopeClient
			} else {
				scopes, err := c.Config.Store.GetAllScopeID()
				if err != nil {
					return errInternal(err.Error())
				}
				if isScopeDiff(td.Scope, scopes) {
					return errUnSupportedScope
				}
			}
		}
	} else if len(td.GrantData.Scope) > 0 {
		td.Scope = td.GrantData.Scope
	} else {
		if len(c.Config.DefaultScopes) == 0 {
			return errNoScope
		}

		td.Scope = c.Config.DefaultScopes
	}

	return nil
}
