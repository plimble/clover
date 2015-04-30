package clover

import (
	"strings"
)

type AuthorizeData struct {
	respType    AuthorizeRespType
	Client      Client
	RedirectURI string
	State       string
	Scope       []string
}

type authorizeCtrl struct {
	clientStore ClientStore
	scopeStore  ScopeStore
	config      *AuthServerConfig
}

func newAuthorizeCtrl(clientStore ClientStore, scopeStore ScopeStore, config *AuthServerConfig) *authorizeCtrl {
	return &authorizeCtrl{clientStore, scopeStore, config}
}

func (a *authorizeCtrl) authorize(ar *authorizeRequest, authRespTypes map[string]AuthorizeRespType, isAuthorized bool, userID string) *Response {
	// the user declined access to the client's application
	if !isAuthorized {
		return errUserDeniedAccess.setRedirect(ar.redirectURI, false, ar.state)
	}

	//validate request
	ad, resp := a.validate(ar, authRespTypes)
	if resp != nil {
		return resp
	}

	return ad.respType.Response(ad, userID)
}

func (a *authorizeCtrl) validate(ar *authorizeRequest, authRespTypes map[string]AuthorizeRespType) (*AuthorizeData, *Response) {
	var resp *Response
	ad := &AuthorizeData{}

	if ad.Client, resp = a.validateClientID(ar); resp != nil {
		return nil, resp
	}

	if ad.RedirectURI, resp = a.validateRedirectURI(ad.Client, ar); resp != nil {
		return nil, resp
	}

	if ad.State, resp = a.validateState(ar); resp != nil {
		return nil, resp
	}

	if ad.respType, resp = a.validateRespType(ad.Client, ar, authRespTypes); resp != nil {
		return nil, resp
	}

	if ad.Scope, resp = a.validateScope(ad.Client, ar); resp != nil {
		resp = resp.setRedirect(ad.RedirectURI, ad.respType.IsImplicit(), ad.State)
		return nil, resp.setRedirect(ad.RedirectURI, ad.respType.IsImplicit(), ad.State)
	}

	return ad, nil
}

func (a *authorizeCtrl) validateClientID(ar *authorizeRequest) (Client, *Response) {
	if ar.clientID == "" {
		return nil, errNoClientID
	}

	//get client
	client, err := a.clientStore.GetClient(ar.clientID)
	if err != nil {
		return nil, errInvalidClientID
	}

	return client, nil
}

func (a *authorizeCtrl) validateState(ar *authorizeRequest) (string, *Response) {
	if a.config.StateParamRequired && ar.state == "" {
		return "", errStateRequired
	}

	return ar.state, nil
}

// Make sure a valid redirect_uri was supplied. If specified, it must match the clientData URI.
func (a *authorizeCtrl) validateRedirectURI(client Client, ar *authorizeRequest) (string, *Response) {
	if ar.redirectURI == "" {
		if client.GetRedirectURI() == "" {
			return "", errNoRedirectURI
		}
		return client.GetRedirectURI(), nil
	}

	if ar.redirectURI != "" && client.GetRedirectURI() != "" && client.GetRedirectURI() != ar.redirectURI {
		return "", errRedirectMismatch
	}

	return ar.redirectURI, nil
}

func (a *authorizeCtrl) validateRespType(client Client, ar *authorizeRequest, authRespTypes map[string]AuthorizeRespType) (AuthorizeRespType, *Response) {
	respType, ok := authRespTypes[ar.responseType]
	if !ok {
		return nil, errInvalidRespType
	}

	if !checkGrantType(client.GetGrantType(), respType.SupportGrant()) {
		return nil, errUnAuthorizedGrant
	}

	return respType, nil
}

func (a *authorizeCtrl) validateScope(client Client, ar *authorizeRequest) ([]string, *Response) {
	scopes := strings.Fields(ar.scope)
	if len(scopes) > 0 {
		if len(client.GetScope()) > 0 {
			if !checkScope(client.GetScope(), scopes...) {
				return nil, errUnSupportedScope
			}
		} else {
			ok, err := a.scopeStore.ExistScopes(scopes...)
			if err != nil {
				return nil, errInternal(err.Error())
			}
			if !ok {
				return nil, errUnSupportedScope
			}
		}
	} else {
		var err error
		scopes, err = a.scopeStore.GetDefaultScope(client.GetClientID())
		if err != nil {
			return nil, errInternal(err.Error())
		}
		if len(scopes) == 0 {
			return nil, errNoScope
		}
	}

	return scopes, nil
}
