package clover

import (
	"net/http"
	"strings"
)

const (
	RESP_TYPE_CODE  = "code"
	RESP_TYPE_TOKEN = "token"
)

type ValidAuthorize func(client Client, scopeIDs []string)

type authorizeRequest struct {
	responseType string
	scope        string
	clientID     string
	redirectURI  string
	state        string
}

func parseAuthorizeRequest(r *http.Request) *authorizeRequest {
	ar := &authorizeRequest{
		state:        r.FormValue("state"),
		redirectURI:  r.FormValue("redirect_uri"),
		responseType: strings.ToLower(r.FormValue("response_type")),
		clientID:     r.FormValue("client_id"),
		scope:        r.FormValue("scope"),
	}

	return ar
}

func (a *AuthorizeServer) Authorize(w http.ResponseWriter, r *http.Request, isAuthorized bool) {
	ar := parseAuthorizeRequest(r)
	resp := a.handleAuthorize(ar, isAuthorized)

	resp.Write(w)
}

func (a *AuthorizeServer) ValidateAuthorize(w http.ResponseWriter, r *http.Request, fn ValidAuthorize) {
	ar := parseAuthorizeRequest(r)

	client, scopes, resp := a.validateAuthorizeRequest(ar)
	if resp != nil {
		resp.Write(w)
		return
	}

	fn(client, scopes)
}

func (a *AuthorizeServer) handleAuthorize(ar *authorizeRequest, isAuthorized bool) *response {
	//re-validate
	client, scopes, resp := a.validateAuthorizeRequest(ar)
	if resp != nil {
		return resp
	}

	// the user declined access to the client's application
	if !isAuthorized {
		return errUserDeniedAccess.SetRedirect(ar.redirectURI, ar.responseType, ar.state)
	}

	return a.RespType[ar.responseType].GetAuthorizeResponse(client, scopes, ar, a)
}

func (a *AuthorizeServer) validateAuthorizeRequest(ar *authorizeRequest) (Client, []string, *response) {
	var resp *response
	var client Client
	var scopes []string

	if client, resp = a.validateAuthorizeClientID(ar); resp != nil {
		return nil, nil, resp
	}

	if resp = a.validateAuthorizeRedirectURI(ar, client); resp != nil {
		return nil, nil, resp
	}

	if resp = a.validateAuthorizeState(ar); resp != nil {
		return nil, nil, resp
	}

	if resp = a.validateAuthorizeResponseType(ar, client); resp != nil {
		resp.SetRedirect(ar.redirectURI, ar.responseType, ar.state)
		return nil, nil, resp
	}

	if scopes, resp = a.validateAuthorizeScope(ar, client); resp != nil {
		resp.SetRedirect(ar.redirectURI, ar.responseType, ar.state)
		return nil, nil, resp
	}

	return client, scopes, nil
}

func (a *AuthorizeServer) validateAuthorizeClientID(ar *authorizeRequest) (Client, *response) {
	if ar.clientID == "" {
		return nil, errNoClientID
	}

	//get client
	client, err := a.Config.Store.GetClient(ar.clientID)
	if err != nil {
		return nil, errInvalidClientID
	}

	return client, nil
}

func (a *AuthorizeServer) validateAuthorizeState(ar *authorizeRequest) *response {
	if a.Config.StateParamRequired && ar.state == "" {
		return errStateRequired
	}

	return nil
}

func (a *AuthorizeServer) validateAuthorizeRedirectURI(ar *authorizeRequest, client Client) *response {
	if ar.redirectURI != "" && client.GetRedirectURI() != "" && client.GetRedirectURI() != ar.redirectURI {
		return errRedirectMismatch
	}

	if ar.redirectURI == "" && client.GetRedirectURI() == "" {
		return errNoRedirectURI
	}

	return nil
}

func (a *AuthorizeServer) validateAuthorizeResponseType(ar *authorizeRequest, client Client) *response {
	switch ar.responseType {
	case "":
		return errInvalidRespType
	case RESP_TYPE_CODE:
		if !checkGrantType(client.GetGrantType(), AUTHORIZATION_CODE) {
			return errUnAuthorizedGrant
		}
		return nil
	case RESP_TYPE_TOKEN:
		if !a.Config.AllowImplicit {
			return errUnSupportedImplicit
		}
		if !checkGrantType(client.GetGrantType(), IMPLICIT) {
			return errUnAuthorizedGrant
		}
		return nil
	}

	return errCodeUnSupportedGrant
}

func (a *AuthorizeServer) validateAuthorizeScope(ar *authorizeRequest, client Client) ([]string, *response) {
	scopes := strings.Fields(ar.scope)

	if len(scopes) > 0 {
		if len(client.GetScope()) == 0 || !checkScope(client.GetScope(), scopes...) {
			return nil, errUnSupportedScope
		}
	} else {
		if len(a.Config.DefaultScopes) == 0 {
			return nil, errNoScope
		}
		return a.Config.DefaultScopes, nil
	}

	return scopes, nil
}
