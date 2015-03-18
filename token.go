package clover

import (
	"encoding/base64"
	"net/http"
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

func (a *AuthorizeServer) Token(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		errMustBePostMetthod.Write(w)
		return
	}

	tr, resp := parseTokenRequest(r, a)
	if resp != nil {
		resp.Write(w)
		return
	}

	resp = a.handleAccessToken(tr)
	resp.Write(w)
}

func parseTokenRequest(r *http.Request, a *AuthorizeServer) (*TokenRequest, *response) {
	clientID, clientSecret, resp := getCredentialsFromHttp(r, a.Config)
	if resp != nil {
		return nil, resp
	}

	tr := &TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        r.FormValue("scope"),
		GrantType:    r.FormValue("grant_type"),
		Username:     r.FormValue("username"),
		Password:     r.FormValue("password"),
		Code:         r.FormValue("code"),
		RedirectURI:  r.FormValue("redirect_uri"),
		RefreshToken: r.FormValue("refresh_token"),
		Assertion:    r.FormValue("assertion"),
	}

	return tr, nil
}

func (a *AuthorizeServer) handleAccessToken(tr *TokenRequest) *response {
	td := &TokenData{}
	if resp := a.validateAccessToken(tr, td); resp != nil {
		return resp
	}

	return td.GrantType.CreateAccessToken(td, a, a.RespType[RESP_TYPE_TOKEN])
}

func (a *AuthorizeServer) validateAccessToken(tr *TokenRequest, td *TokenData) *response {
	var resp *response
	if td.GrantType, resp = a.validateAccessTokenGrantType(tr); resp != nil {
		return resp
	}

	if td.GrantData, resp = td.GrantType.Validate(tr, a); resp != nil {
		return resp
	}

	if td.GrantType.GetGrantType() != CLIENT_CREDENTIALS {
		if resp = a.validateAccessTokenClient(tr, td.GrantData); resp != nil {
			return resp
		}
	}

	if !checkGrantType(td.GrantData.GrantType, td.GrantType.GetGrantType()) {
		return errUnAuthorizedGrant
	}

	if td.Scope, resp = a.validateAccessTokenScope(tr, td.GrantData); resp != nil {
		return resp
	}

	return nil
}

func (a *AuthorizeServer) validateAccessTokenGrantType(tr *TokenRequest) (GrantType, *response) {
	if tr.GrantType == "" {
		return nil, errGrantTypeRequired
	}

	if _, ok := a.Grant[tr.GrantType]; !ok {
		return nil, errUnSupportedGrantType
	}

	return a.Grant[tr.GrantType], nil
}

func (a *AuthorizeServer) validateAccessTokenClient(tr *TokenRequest, grantData *GrantData) *response {
	var err error
	client, err := a.Config.Store.GetClient(tr.ClientID)
	if err != nil {
		return errInternal(err.Error())
	}

	if grantData.ClientID != "" && grantData.ClientID != client.GetClientID() {
		return errInvalidClientCredentail
	}

	copyClientToGrant(grantData, client)

	return nil
}

func (a *AuthorizeServer) validateAccessTokenScope(tr *TokenRequest, grantData *GrantData) ([]string, *response) {
	scopes := strings.Fields(tr.Scope)
	if len(scopes) > 0 {
		if len(grantData.Scope) > 0 {
			if !checkScope(grantData.Scope, scopes...) {
				return nil, errInvalidScopeRequest
			}
		} else {
			return nil, errNoScope
		}
	} else if len(grantData.Scope) > 0 {
		scopes = grantData.Scope
	} else {
		if len(a.Config.DefaultScopes) == 0 {
			return nil, errNoScope
		}

		scopes = a.Config.DefaultScopes
	}

	return scopes, nil
}

func copyClientToGrant(grantData *GrantData, client Client) {
	if grantData.UserID == "" {
		grantData.UserID = client.GetUserID()
	}

	if grantData.ClientID == "" {
		grantData.ClientID = client.GetClientID()
	}

	if len(grantData.Scope) == 0 {
		grantData.Scope = client.GetScope()
	}

	if len(grantData.GrantType) == 0 {
		grantData.GrantType = client.GetGrantType()
	}
}

func getCredentialsFromHttp(r *http.Request, config *Config) (string, string, *response) {
	headerAuth := r.Header.Get("Authorization")

	switch {
	case headerAuth != "":
		s := strings.SplitN(headerAuth, " ", 2)
		if len(s) != 2 || s[0] != "Basic" {
			return "", "", errInvalidAuthHeader
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			return "", "", errInternal(err.Error())
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			return "", "", errInvalidAuthMSG
		}

		return pair[0], pair[1], nil
	case config.AllowCredentialsBody:
		if r.PostForm.Get(`client_id`) == "" || r.PostForm.Get(`client_secret`) == "" {
			return "", "", errCredentailsNotInBody
		}

		return r.PostForm.Get("client_id"), r.PostForm.Get("client_secret"), nil
	}

	return "", "", errCredentailsRequired
}
