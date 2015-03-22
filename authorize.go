package clover

import (
	"strings"
)

type authController struct {
	config        *AuthConfig
	authRespType  AuthRespType
	tokenRespType AuthRespType
}

type authorizeRequest struct {
	responseType string
	scope        string
	clientID     string
	redirectURI  string
	state        string
}

func newAuthController(config *AuthConfig, authRespType, tokenRespType AuthRespType) *authController {
	return &authController{config, authRespType, tokenRespType}
}

func (a *authController) handleAuthorize(ar *authorizeRequest, isAuthorized bool) *Response {
	//re-validate
	client, scopes, respType, resp := a.validateAuthorizeRequest(ar)
	if resp != nil {
		return resp
	}

	// the user declined access to the client's application
	if !isAuthorized {
		return errUserDeniedAccess.setRedirect(ar.redirectURI, ar.responseType, ar.state)
	}

	return respType.GetAuthResponse(ar, client, scopes)
}

func (a *authController) validateAuthorizeRequest(ar *authorizeRequest) (Client, []string, AuthRespType, *Response) {
	var resp *Response
	var client Client
	var scopes []string
	var respType AuthRespType

	if client, resp = a.validateAuthorizeClientID(ar); resp != nil {
		return nil, nil, nil, resp
	}

	if resp = a.validateAuthorizeRedirectURI(ar, client); resp != nil {
		return nil, nil, nil, resp
	}

	if resp = a.validateAuthorizeState(ar); resp != nil {
		return nil, nil, nil, resp
	}

	if respType, resp = a.validateAuthorizeResponseType(ar, client); resp != nil {
		resp.setRedirect(ar.redirectURI, ar.responseType, ar.state)
		return nil, nil, nil, resp
	}

	if scopes, resp = a.validateAuthorizeScope(ar, client); resp != nil {
		resp.setRedirect(ar.redirectURI, ar.responseType, ar.state)
		return nil, nil, nil, resp
	}

	return client, scopes, respType, nil
}

func (a *authController) validateAuthorizeClientID(ar *authorizeRequest) (Client, *Response) {
	if ar.clientID == "" {
		return nil, errNoClientID
	}

	//get client
	client, err := a.config.AuthServerStore.GetClient(ar.clientID)
	if err != nil {
		return nil, errInvalidClientID
	}

	return client, nil
}

func (a *authController) validateAuthorizeState(ar *authorizeRequest) *Response {
	if a.config.StateParamRequired && ar.state == "" {
		return errStateRequired
	}

	return nil
}

// Make sure a valid redirect_uri was supplied. If specified, it must match the clientData URI.
func (a *authController) validateAuthorizeRedirectURI(ar *authorizeRequest, client Client) *Response {
	if ar.redirectURI == "" {
		if client.GetRedirectURI() == "" {
			return errNoRedirectURI
		}
		ar.redirectURI = client.GetRedirectURI()
		return nil
	}

	if ar.redirectURI != "" && client.GetRedirectURI() != "" && client.GetRedirectURI() != ar.redirectURI {
		return errRedirectMismatch
	}

	return nil
}

func (a *authController) validateAuthorizeResponseType(ar *authorizeRequest, client Client) (AuthRespType, *Response) {
	switch ar.responseType {
	case "":
		return nil, errInvalidRespType
	case RESP_TYPE_CODE:
		if !checkGrantType(client.GetGrantType(), AUTHORIZATION_CODE) {
			return nil, errUnAuthorizedGrant
		}
		return a.authRespType, nil
	case RESP_TYPE_TOKEN:
		if !a.config.AllowImplicit {
			return nil, errUnSupportedImplicit
		}

		if !checkGrantType(client.GetGrantType(), IMPLICIT) {
			return nil, errUnAuthorizedGrant
		}
		return a.tokenRespType, nil
	}

	return nil, errCodeUnSupportedGrant
}

func (a *authController) validateAuthorizeScope(ar *authorizeRequest, client Client) ([]string, *Response) {
	scopes := strings.Fields(ar.scope)

	if len(scopes) > 0 {
		if len(client.GetScope()) == 0 || !checkScope(client.GetScope(), scopes...) {
			return nil, errUnSupportedScope
		}
	} else {
		if len(a.config.DefaultScopes) == 0 {
			return nil, errNoScope
		}
		return a.config.DefaultScopes, nil
	}

	return scopes, nil
}
