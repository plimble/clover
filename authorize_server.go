package clover

import (
	"net/http"
	"strings"
)

const (
	RESP_TYPE_CODE  = "code"
	RESP_TYPE_TOKEN = "token"
)

type authorizeRequest struct {
	responseType string
	scope        string
	clientID     string
	redirectURI  string
	state        string
}

func (a *authorizeRequest) GetResponseType() string {
	return a.responseType
}

func (a *authorizeRequest) GetState() string {
	return a.state
}

func (a *authorizeRequest) GetRedirectURI() string {
	return a.redirectURI
}

type AuthorizeData struct {
	Code         string
	Client       *Client
	ResponseType string
	State        string
	RedirectURI  string
	Scope        []string
	TokenData
}

func (a *AuthorizeData) GetResponseType() string {
	return a.ResponseType
}

func (a *AuthorizeData) GetState() string {
	return a.State
}

func (a *AuthorizeData) GetRedirectURI() string {
	return a.RedirectURI
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

func (c *Clover) Authorize(w http.ResponseWriter, r *http.Request, isAuthorized bool) {
	ar := parseAuthorizeRequest(r)
	resp := c.handleAuthorize(ar, isAuthorized)

	resp.Write(w)
}

func (c *Clover) ValidateAuthorize(w http.ResponseWriter, r *http.Request) *AuthorizeData {
	ar := parseAuthorizeRequest(r)

	ad, resp := c.validateAuthorizeRequest(ar)
	if resp != nil {
		resp.Write(w)
		return nil
	}

	return ad
}

func (c *Clover) handleAuthorize(ar *authorizeRequest, isAuthorized bool) *Response {
	//re-validate
	ad, resp := c.validateAuthorizeRequest(ar)
	if resp != nil {
		return resp
	}

	// the user declined access to the client's application
	if !isAuthorized {
		return errUserDeniedAccess.SetRedirect(ar)
	}

	return c.RespType[ad.ResponseType].GetAuthorizeResponse(ad, c)
}

func (c *Clover) validateAuthorizeRequest(ar *authorizeRequest) (*AuthorizeData, *Response) {
	ad := &AuthorizeData{}
	if resp := c.validateAuthorizeClientID(ar, ad); resp != nil {
		return nil, resp
	}

	if resp := c.validateAuthorizeRedirectURI(ar, ad); resp != nil {
		return nil, resp
	}

	if resp := c.validateAuthorizeState(ar, ad); resp != nil {
		return nil, resp
	}

	if resp := c.validateAuthorizeResponseType(ar, ad); resp != nil {
		return nil, resp
	}

	if resp := c.validateAuthorizeScope(ar, ad); resp != nil {
		return nil, resp
	}

	return ad, nil
}

func (c *Clover) validateAuthorizeClientID(ar *authorizeRequest, ad *AuthorizeData) *Response {
	if ar.clientID == "" {
		return errNoClientID
	}

	//get client
	client, err := c.Config.Store.GetClient(ar.clientID)
	if err != nil {
		return errInvalidClientID
	}

	ad.Client = client
	return nil
}

func (c *Clover) validateAuthorizeState(ar *authorizeRequest, ad *AuthorizeData) *Response {
	if c.Config.StateParamRequired && ar.state == "" {
		return errStateRequired
	}

	ad.State = ar.state
	return nil
}

func (c *Clover) validateAuthorizeRedirectURI(ar *authorizeRequest, ad *AuthorizeData) *Response {
	ad.RedirectURI = ar.redirectURI

	if ar.redirectURI != "" && ad.Client.RedirectURI != "" && ad.Client.RedirectURI != ar.redirectURI {
		return errRedirectMismatch
	}

	if ar.redirectURI == "" && ad.Client.RedirectURI == "" {
		return errNoRedirectURI
	}

	ad.RedirectURI = ad.Client.RedirectURI
	return nil
}

func (c *Clover) validateAuthorizeResponseType(ar *authorizeRequest, ad *AuthorizeData) *Response {
	ad.ResponseType = ar.responseType
	switch ad.ResponseType {
	case "":
		return errInvalidRespType.SetRedirect(ad)
	case "code":
		if !checkGrantType(ad.Client.GrantType, AUTHORIZATION_CODE) {
			return errUnAuthorizedGrant.SetRedirect(ad)
		}
		return nil
	case "token":
		if !c.Config.AllowImplicit {
			return errUnSupportedImplicit.SetRedirect(ad)
		}
		if !checkGrantType(ad.Client.GrantType, IMPLICIT) {
			return errUnAuthorizedGrant.SetRedirect(ad)
		}
		return nil
	}

	return errCodeUnSupportedGrant.SetRedirect(ad)
}

func (c *Clover) validateAuthorizeScope(ar *authorizeRequest, ad *AuthorizeData) *Response {
	ad.Scope = strings.Fields(ar.scope)

	if len(ad.Scope) > 0 {
		if len(ad.Client.Scope) == 0 || !checkScope(ad.Scope, ad.Client.Scope) {
			return errUnSupportedScope.SetRedirect(ad)
		}
	} else {
		if len(c.Config.DefaultScopes) == 0 {
			return errNoScope.SetRedirect(ad)
		}
		ad.Scope = c.Config.DefaultScopes
		return nil
	}

	return nil
}
